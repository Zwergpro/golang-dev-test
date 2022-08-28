package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"homework-1/config"
	"homework-1/internal/api/storage"
	"homework-1/internal/metrics"
	postgresRepository "homework-1/internal/repository/postgres"
	pbStorage "homework-1/pkg/api/storage/v1"
	"net"
	"net/http"
	"os"
)

// сервис для работы с базами данных
func main() {
	SetUpLogger()

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
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	if err = pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	poolConfig := pool.Config()
	poolConfig.MaxConnIdleTime = config.DBMaxConnIdleTime
	poolConfig.MaxConnLifetime = config.DBMaxConnLifetime
	poolConfig.MinConns = config.DBMinConns
	poolConfig.MaxConns = config.DBMaxConns

	grpcServer := grpc.NewServer()

	appMetrics := metrics.NewMetrics()
	appMetrics.Publish()

	go func() {
		log.Infof("starting metrics http server on %s", config.StorageStatAddress)
		if err = http.ListenAndServe(config.StorageStatAddress, nil); err != nil {
			log.Fatal(err)
		}
	}()

	deps := storage.Deps{
		ProductRepository: postgresRepository.NewRepository(pool),
		Metrics:           appMetrics,
	}

	pbStorage.RegisterStorageServiceServer(grpcServer, storage.New(deps))

	listener, err := net.Listen("tcp", config.StorageServiceAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("starting grpc server on %s", config.StorageServiceAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
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
