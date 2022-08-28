package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	"homework-1/internal/api/proxyApi"
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

	err := opentelemetry.SetGlobalTracer("proxy-apy", config.TracerUrl)
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

	deps := proxyApi.Deps{
		StorageClient: client,
		Metrics:       appMetrics,
	}
	pbApi.RegisterApiServiceServer(grpcServer, proxyApi.New(deps))

	listener, err := net.Listen("tcp", config.ProxyApiServiceAddress)
	if err != nil {
		log.WithError(err).Fatal("failed to listen proxy api service")
	}
	log.Infof("starting grpc server on %s", config.ProxyApiServiceAddress)
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
