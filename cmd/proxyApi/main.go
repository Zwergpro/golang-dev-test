package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"homework-1/config"
	"homework-1/internal/api/proxyApi"
	pbStorage "homework-1/pkg/api/storage/v1"
	pbApi "homework-1/pkg/api/v1"
	"log"
	"net"
)

func main() {
	grpcServer := grpc.NewServer()

	conn, err := grpc.Dial(config.StorageServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	client := pbStorage.NewStorageServiceClient(conn)

	deps := proxyApi.Deps{
		StorageClient: client,
	}
	pbApi.RegisterApiServiceServer(grpcServer, proxyApi.New(deps))

	listener, err := net.Listen("tcp", config.ProxyApiServiceAddress)
	if err != nil {
		log.Fatal(err)
	}
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
