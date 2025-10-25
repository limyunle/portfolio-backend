package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/limyunle/portfolio-backend/internal/aggregate"
	"github.com/limyunle/portfolio-backend/internal/config"
	"github.com/limyunle/portfolio-backend/internal/s3"
)

var (
	s3Service  s3.Storage
	bucketName string
)

func init() {
	cfg := config.LoadConfig()
	bucketName = cfg.S3Bucket
	s3Service = s3.NewService(cfg.S3Service, bucketName)

	if bucketName == "" {
		log.Fatal("S3_BUCKET environment variable is not set")
	}
	log.Println("Backend Lambda initialized")
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	start := time.Now()
	log.Println("API request received")

	if req.Path != "/aggregate/stats" {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusNotFound,
			Body:       "Not Found",
		}, nil
	}

	var stats aggregate.CombinedStats
	if err := s3Service.GetJSON(ctx, bucketName, "aggregate-stats.json", &stats); err != nil {
		log.Printf("Failed to fetch from S3: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error fetching data: %v", err),
		}, nil
	}

	respBody, err := json.Marshal(stats)
	if err != nil {
		log.Printf("Failed to marshal response: %v", err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       fmt.Sprintf("Error encoding response: %v", err),
		}, nil
	}

	log.Printf("Request served in %s", time.Since(start))
	return events.APIGatewayProxyResponse{
		StatusCode:      http.StatusOK,
		Body:            string(respBody),
		Headers:         map[string]string{"Content-Type": "application/json"},
		IsBase64Encoded: false,
	}, nil
}

func main() {
	lambda.Start(handler)
}
