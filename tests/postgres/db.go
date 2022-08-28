//go:build integration
// +build integration

package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strings"
	"sync"
	"testing"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	"homework-1/tests/config"
	"time"
)

type TDB struct {
	sync.Mutex
	Pool *pgxpool.Pool
}

func NewFromEnv(ctx context.Context) *TDB {
	cfg, err := config.FromEnv()
	if err != nil {
		log.Fatal(err)
	}

	psqlConn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName)

	pool, err := pgxpool.Connect(ctx, psqlConn)
	if err != nil {
		log.Fatal("can't connect to database", err)
	}

	if err = pool.Ping(ctx); err != nil {
		log.Fatal("ping database error", err)
	}

	poolConfig := pool.Config()
	poolConfig.MaxConnIdleTime = time.Minute
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MinConns = 2
	poolConfig.MaxConns = 6

	return &TDB{Pool: pool}
}

func (d *TDB) SetUp(t *testing.T) {
	t.Helper()
	ctx := context.Background()
	d.Lock()
	d.Truncate(ctx)
}

func (d *TDB) TearDown() {
	defer d.Unlock()
	d.Truncate(context.Background())
}

func (d *TDB) Truncate(ctx context.Context) {
	var tables []string

	tablesQuery := "SELECT table_name FROM information_schema.tables WHERE table_schema='public' AND table_type='BASE TABLE' AND table_name != 'goose_db_version'"
	if err := pgxscan.Select(ctx, d.Pool, &tables, tablesQuery); err != nil {
		panic(err)
	}

	if len(tables) == 0 {
		log.Fatal("run migration plz")
	}
	q := fmt.Sprintf("Truncate table %s", strings.Join(tables, ","))
	if _, err := d.Pool.Exec(ctx, q); err != nil {
		panic(err)
	}
}
