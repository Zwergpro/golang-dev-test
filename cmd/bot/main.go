package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework-1/internal/commander"
	"homework-1/internal/handlers"
	postgresRepository "homework-1/internal/repository/postgres"
	"log"
	"os"
)

func main() {
	tgApiKey := os.Getenv("TG_API_KEY")
	if tgApiKey == "" {
		log.Fatal("TG_API_KEY env variable does not exist")
	}

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

	cmd, err := commander.Init(tgApiKey, postgresRepository.NewRepository(pool))
	if err != nil {
		log.Fatal(err)
	}

	handlers.AddHandlers(cmd)

	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
