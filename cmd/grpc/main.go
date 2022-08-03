package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"homework-1/config"
	"homework-1/internal/api"
	postgresRepository "homework-1/internal/repository/postgres"
	pb "homework-1/pkg/api/v1"
	"log"
	"net"
)

func main() {
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

	deps := api.Deps{
		ProductRepository: postgresRepository.NewRepository(pool),
	}

	pb.RegisterAdminServiceServer(grpcServer, api.New(deps))

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	if err = grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
