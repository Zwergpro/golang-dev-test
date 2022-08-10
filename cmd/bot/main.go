package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"homework-1/config"
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

	cmd, err := commander.Init(tgApiKey, postgresRepository.NewRepository(pool))
	if err != nil {
		log.Fatal(err)
	}

	handlers.AddHandlers(cmd)

	if err = cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
