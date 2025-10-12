package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/limyunle/portfolio-backend/internal/aggregate"
	"github.com/limyunle/portfolio-backend/internal/config"
	"github.com/limyunle/portfolio-backend/internal/github"
	"github.com/limyunle/portfolio-backend/internal/leetcode"
	"github.com/limyunle/portfolio-backend/internal/s3"
)

func LambdaHandler(ctx context.Context) error {
	start := time.Now()
	log.Println("Lambda invoked")

	cfg := config.LoadConfig()
	log.Printf("Loaded config: Username=%s, S3Bucket=%s\n", cfg.Username, cfg.S3Bucket)

	if cfg.Username == "" {
		return fmt.Errorf("environment variable USERNAME is not set")
	}
	if cfg.S3Bucket == "" {
		return fmt.Errorf("environment variable S3_BUCKET is not set")
	}

	log.Println("Initializing GitHub service...")
	githubService := github.NewService()
	log.Println("Initializing LeetCode service...")
	leetcodeService := leetcode.NewService()
	log.Println("Initializing S3 service...")
	s3Service := s3.NewService(cfg.S3Service, cfg.S3Bucket)

	aggregateService := &aggregate.Service{
		GitHubService:   githubService,
		LeetCodeService: leetcodeService,
		S3Service:       s3Service,
		BucketName:      cfg.S3Bucket,
	}

	if err := aggregateService.RefreshAndStore(cfg.Username); err != nil {
		return fmt.Errorf("failed to refresh aggregated stats: %w", err)
	}

	log.Println("Aggregated stats successfully refreshed and uploaded to S3")
	log.Printf("Lambda finished in %s\n", time.Since(start))

	return nil
}

func main() {
	lambda.Start(LambdaHandler)
}
