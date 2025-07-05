package config

import (
	"log"

	"github.com/caarlos0/env/v10"
)

type Config struct {
	MINIO_ENDPOINT    string `env:"MINIO_ENDPOINT" env-required:"true" env-default:"localhost:9000"`
	MINIO_ACCESS_KEY  string `env:"MINIO_ACCESS_KEY" env-required:"true"`
	MINIO_SECRET_KEY  string `env:"MINIO_SECRET_KEY" env-required:"true"`
	MINIO_BUCKET_NAME string `env:"MINIO_BUCKET_NAME" env-required:"true"`
}

func MustHaveConfig() Config {
	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	return cfg
}
