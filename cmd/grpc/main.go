package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"homework-1/internal/api"
	postgresRepository "homework-1/internal/repository/postgres"
	pb "homework-1/pkg/api/v1"
	"log"
	"net"
	"time"
)

const (
	// Database config
	Host     = "localhost"
	Port     = 6432
	User     = "postgres"
	Password = "postgres"
	DBname   = "postgres"

	MaxConnIdleTime = time.Minute
	MaxConnLifetime = time.Hour
	MinConns        = 2
	MaxConns        = 4
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Host, Port, User, Password, DBname)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	config := pool.Config()
	config.MaxConnIdleTime = MaxConnIdleTime
	config.MaxConnLifetime = MaxConnLifetime
	config.MinConns = MinConns
	config.MaxConns = MaxConns

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
