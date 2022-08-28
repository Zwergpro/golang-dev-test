package main

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	"homework-1/internal/api/proxyApi"
	"homework-1/internal/metrics"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"net"
	"net/http"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	grpcServer := grpc.NewServer()

	conn, err := grpc.Dial(config.StorageServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	deps := proxyApi.Deps{
		StorageClient: client,
		Metrics:       appMetrics,
	}
	pbApi.RegisterApiServiceServer(grpcServer, proxyApi.New(deps))

	listener, err := net.Listen("tcp", config.ProxyApiServiceAddress)
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("starting grpc server on %s", config.ProxyApiServiceAddress)
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
