package config

import (
	"context"
	"log"

	"github.com/SidM81/goShare/models"
	"github.com/caarlos0/env/v10"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Config struct {
	SERVER_ADDRESS    string `env:"SERVER_ADDRESS" env-required:"true" env-default:"8080"`
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

var DB *gorm.DB

func InitDatabase() {
	db, err := gorm.Open(sqlite.Open("storage/goshare.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to SQLite: %v", err)
	}

	DB = db
	db.AutoMigrate(&models.File{}, &models.Chunk{})
}

func GetDB() *gorm.DB {
	return DB
}

var MinioClient *minio.Client

func InitMinio(cfg Config) {
	client, err := minio.New(cfg.MINIO_ENDPOINT, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.MINIO_ACCESS_KEY, cfg.MINIO_SECRET_KEY, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
	}

	MinioClient = client

	exists, err := client.BucketExists(context.Background(), cfg.MINIO_BUCKET_NAME)
	if err != nil {
		log.Fatalf("Error checking bucket: %v", err)
	}
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.MINIO_BUCKET_NAME, minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Error creating bucket: %v", err)
		}
	}
}

func GetMinioClient() *minio.Client {
	return MinioClient
}
