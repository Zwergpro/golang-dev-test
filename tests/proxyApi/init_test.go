//go:build integration
// +build integration

package proxyApi

import (
	"context"
	pbApi "homework-1/pkg/api/v1"
	"log"
	"time"

	"google.golang.org/grpc"
	"homework-1/tests/config"
	"homework-1/tests/postgres"
)

var (
	ProxyApiClient pbApi.ApiServiceClient
	DB             *postgres.TDB
)

func init() {
	cfg, err := config.FromEnv()

	conn, err := grpc.Dial(cfg.Host, grpc.WithInsecure(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	ProxyApiClient = pbApi.NewApiServiceClient(conn)

	DB = postgres.NewFromEnv(context.Background())
}
