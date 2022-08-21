//go:build integration
// +build integration

package storage

import (
	"context"
	pbStorage "homework-1/pkg/api/storage/v1"
	"log"
	"time"

	"google.golang.org/grpc"
	"homework-1/tests/postgres"
)

var (
	StorageClient pbStorage.StorageServiceClient
	DB            *postgres.TDB
)

func init() {
	//cfg, err := config.FromEnv()

	conn, err := grpc.Dial("0.0.0.0:8080", grpc.WithInsecure(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	StorageClient = pbStorage.NewStorageServiceClient(conn)

	DB = postgres.NewFromEnv(context.Background())
}
