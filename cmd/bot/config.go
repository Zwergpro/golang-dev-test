package main

import "time"

const (
	Host     = "0.0.0.0"
	Port     = 6432
	User     = "postgres"
	Password = "postgres"
	DBname   = "postgres"

	MaxConnIdleTime = time.Minute
	MaxConnLifetime = time.Hour
	MinConns        = 2
	MaxConns        = 6
)
