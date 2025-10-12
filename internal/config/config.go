package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port             int
	S3Bucket         string
	S3Service        *s3.Client
	RefreshFrequency int
}

func LoadConfig() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	portStr := os.Getenv("PORT")
	if portStr == "" {
		portStr = "8080"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid PORT value: %v", err)
	}

	bucket := os.Getenv("PORTFOLIO_S3_BUCKET")
	if bucket == "" {
		log.Fatal("PORTFOLIO_S3_BUCKET not set")
	}

	awsCfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("unable to load AWS config: %v", err)
	}

	s3Client := s3.NewFromConfig(awsCfg)

	return &AppConfig{
		Port:      port,
		S3Bucket:  bucket,
		S3Service: s3Client,
	}
}
