//go:build integration
// +build integration

package config

import "github.com/kelseyhightower/envconfig"

const envPrefix = "QA"

type Config struct {
	ProxyApiHost string `split_words:"true" default:":8081"`
	StorageHost  string `split_words:"true" default:":8080"`
	DBHost       string `split_words:"true" default:"localhost"`
	DBPort       int    `split_words:"true" default:"6432"`
	DBUser       string `split_words:"true" default:"postgres"`
	DBPassword   string `split_words:"true" default:"postgres"`
	DBName       string `split_words:"true" default:"postgres"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}
