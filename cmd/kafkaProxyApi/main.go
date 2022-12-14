package main

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	"homework-1/internal/api/kafkaProxyApi"
	redisCache "homework-1/internal/cache/redis"
	"homework-1/internal/metrics"
	"homework-1/internal/opentelemetry"
	pbStorage "homework-1/pkg/api/storage/v2"
	pbApi "homework-1/pkg/api/v2"
	"net"
	"net/http"
	"os"
)

func main() {
	SetUpLogger()

	err := opentelemetry.SetGlobalTracer("kafka-proxy-apy", config.TracerUrl)
	if err != nil {
		log.WithError(err).Fatal("failed to create tracer")
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(opentelemetry.UnaryServerInterceptor()),
		grpc.StreamInterceptor(opentelemetry.StreamServerInterceptor()),
	)

	conn, err := grpc.Dial(
		config.StorageServiceAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(opentelemetry.UnaryClientInterceptor()),
		grpc.WithStreamInterceptor(opentelemetry.StreamClientInterceptor()),
	)
	if err != nil {
		log.WithError(err).Fatal("failed to connect to storage service")
	}

	client := pbStorage.NewStorageServiceClient(conn)

	appMetrics := metrics.NewMetrics()
	appMetrics.Publish()

	go func() {
		log.Infof("starting metrics http server on %s", config.ProxyApiStatAddress)
		if err = http.ListenAndServe(config.ProxyApiStatAddress, nil); err != nil {
			log.WithError(err).Fatal("failed to start metrics http server")
		}
	}()

	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	syncProducer, err := sarama.NewSyncProducer(config.GetKafkaBrokers(), cfg)
	if err != nil {
		log.WithError(err).Fatal("kafka: NewSyncProducer")
	}
	syncProducer = otelsarama.WrapSyncProducer(cfg, syncProducer)
	cache := redisCache.New(config.GetRedisOpts(), appMetrics)

	deps := kafkaProxyApi.Deps{
		StorageClient: client,
		Metrics:       appMetrics,
		Producer:      syncProducer,
		Cache:         cache,
	}
	pbApi.RegisterApiServiceServer(grpcServer, kafkaProxyApi.New(deps))

	listener, err := net.Listen("tcp", config.ProxyApiServiceAddress)
	if err != nil {
		log.WithError(err).Fatal("failed to listen")
	}
	log.Infof("starting grpc server on %s", config.ProxyApiServiceAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.WithError(err).Fatal("failed to serve")
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
