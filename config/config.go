package config

import "time"

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
)
