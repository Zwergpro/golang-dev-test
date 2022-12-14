package config

import (
	"github.com/go-redis/redis/v9"
	"time"
)

const (
	DBHost     = "0.0.0.0"
	DBPort     = 6432
	DBUser     = "postgres"
	DBPassword = "postgres"
	DBName     = "postgres"

	DBMaxConnIdleTime = time.Minute
	DBMaxConnLifetime = time.Hour
	DBMinConns        = 2
	DBMaxConns        = 6
)

const (
	StorageServiceAddress     = ":8080"
	StorageStatAddress        = ":9080"
	ProxyApiServiceAddress    = ":8081"
	ProxyApiStatAddress       = ":9081"
	HTTPGatewayServiceAddress = ":8082"

	TracerUrl = "http://localhost:14268/api/traces"
)

const (
	RedisAddr = "localhost:6379"
	RedisDB   = 0
	RedisPass = ""
)

func GetKafkaBrokers() []string {
	return []string{"localhost:29091", "localhost:19091", "localhost:39091"}
}

func GetRedisOpts() *redis.Options {
	return &redis.Options{Addr: RedisAddr, DB: RedisDB, Password: RedisPass}
}
