package main

import (
	"google.golang.org/grpc"
	"homework-1/internal/api"
	pb "homework-1/pkg/api/v1"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAdminServiceServer(grpcServer, api.New())

	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
