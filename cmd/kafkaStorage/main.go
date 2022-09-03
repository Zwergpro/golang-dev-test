package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"homework-1/config"
	"homework-1/internal/api/kafkaStorage"
	"homework-1/internal/api/kafkaStorage/consumers"
	redisCache "homework-1/internal/cache/redis"
	"homework-1/internal/metrics"
	"homework-1/internal/opentelemetry"
	"homework-1/internal/repository"
	postgresRepository "homework-1/internal/repository/postgres"
	pbStorage "homework-1/pkg/api/storage/v2"
	"net"
	"net/http"
	"os"
)

// сервис для работы с базами данных
func main() {
	SetUpLogger()

	err := opentelemetry.SetGlobalTracer("kafka-storage", config.TracerUrl)
	if err != nil {
		log.WithError(err).Fatal("failed to create tracer")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to postgres")
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.WithError(err).Fatal("failed to ping postgres")
	}

	poolConfig := pool.Config()
	poolConfig.MaxConnIdleTime = config.DBMaxConnIdleTime
	poolConfig.MaxConnLifetime = config.DBMaxConnLifetime
	poolConfig.MinConns = config.DBMinConns
	poolConfig.MaxConns = config.DBMaxConns

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(opentelemetry.UnaryServerInterceptor()),
		grpc.StreamInterceptor(opentelemetry.StreamServerInterceptor()),
	)

	appMetrics := metrics.NewMetrics()
	appMetrics.Publish()

	go func() {
		log.Infof("starting metrics http server on %s", config.StorageStatAddress)
		if err = http.ListenAndServe(config.StorageStatAddress, nil); err != nil {
			log.WithError(err).Fatal("failed to start metrics http server")
		}
	}()

	runStorageKafkaConsumers(postgresRepository.NewRepository(pool), appMetrics)

	cache := redisCache.New(&redis.Options{Addr: "localhost:6379", DB: 0, Password: ""}, appMetrics)

	deps := kafkaStorage.Deps{
		ProductRepository: postgresRepository.NewRepository(pool),
		Metrics:           appMetrics,
		Cache:             cache,
	}

	pbStorage.RegisterStorageServiceServer(grpcServer, kafkaStorage.New(deps))

	listener, err := net.Listen("tcp", config.StorageServiceAddress)
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}
	log.Infof("starting grpc server on %s", config.StorageServiceAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.WithError(err).Fatal("failed to start grpc server")
	}
}

func SetUpLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	if os.Getenv("QA_DEBUG") == "True" {
		log.SetLevel(log.DebugLevel)
	}
}

func runStorageKafkaConsumers(productRepository repository.Product, appMetrics *metrics.Metrics) {
	productCreateConsumer := &consumers.ProductCreateConsumer{
		ProductRepository: productRepository,
		Metrics:           appMetrics,
	}
	go productCreateConsumer.StartConsuming(context.Background())

	productUpdateConsumer := &consumers.ProductUpdateConsumer{
		ProductRepository: productRepository,
		Metrics:           appMetrics,
	}
	go productUpdateConsumer.StartConsuming(context.Background())

	productDeleteConsumer := &consumers.ProductDeleteConsumer{
		ProductRepository: productRepository,
		Metrics:           appMetrics,
	}
	go productDeleteConsumer.StartConsuming(context.Background())
}
