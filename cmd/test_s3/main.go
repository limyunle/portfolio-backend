package main

import (
	"context"
	"fmt"
	"log"

	"github.com/limyunle/portfolio-backend/internal/config"
	"github.com/limyunle/portfolio-backend/internal/s3"
)

func main() {
	appCfg := config.LoadConfig()

	s3Service := s3.NewService(appCfg.S3Service, appCfg.S3Bucket)

	ctx := context.Background()

	type Profile struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	testData := Profile{Name: "Yun Le Lim", Email: "limyunle@gmail.com"}

	key := "test/profile.json"
	if err := s3Service.UploadJSON(ctx, appCfg.S3Bucket, key, testData); err != nil {
		log.Fatalf("upload failed: %v", err)
	}
	fmt.Println("Successfully uploaded:", key)

	var downloaded Profile
	if err := s3Service.GetJSON(ctx, appCfg.S3Bucket, key, &downloaded); err != nil {
		log.Fatalf("download failed: %v", err)
	}

	fmt.Println("Retrieved data:", downloaded)
}
