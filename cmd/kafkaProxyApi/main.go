package main

import (
	"github.com/Shopify/sarama"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/Shopify/sarama/otelsarama"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	"homework-1/internal/api/kafkaProxyApi"
	"homework-1/internal/metrics"
	"homework-1/internal/opentelemetry"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"net"
	"net/http"
	"os"
)

func main() {
	SetUpLogger()

	err := opentelemetry.SetGlobalTracer("kafka-proxy-apy", "http://localhost:14268/api/traces")
	if err != nil {
		log.Fatalf("failed to create tracer: %v", err)
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
		log.Fatal(err)
	}

	client := pbStorage.NewStorageServiceClient(conn)

	appMetrics := metrics.NewMetrics()
	appMetrics.Publish()

	go func() {
		log.Infof("starting metrics http server on %s", config.ProxyApiStatAddress)
		if err = http.ListenAndServe(config.ProxyApiStatAddress, nil); err != nil {
			log.Fatal(err)
		}
	}()

	brokers := []string{"localhost:29091", "localhost:19091", "localhost:39091"}
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(brokers, cfg)
	if err != nil {
		log.Fatalf("sync kafka: %v", err)
	}

	syncProducer = otelsarama.WrapSyncProducer(cfg, syncProducer)

	deps := kafkaProxyApi.Deps{
		StorageClient: client,
		Metrics:       appMetrics,
		Producer:      syncProducer,
	}
	pbApi.RegisterApiServiceServer(grpcServer, kafkaProxyApi.New(deps))

	listener, err := net.Listen("tcp", config.ProxyApiServiceAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("starting grpc server on %s", config.ProxyApiServiceAddress)
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
