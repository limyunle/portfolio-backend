package main

import (
	"log"
	"strconv"

	"github.com/limyunle/portfolio-backend/internal/config"
	"github.com/limyunle/portfolio-backend/internal/routes"
	"github.com/limyunle/portfolio-backend/internal/s3"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	appCfg := config.LoadConfig()

	r := gin.Default()
	r.Use(cors.Default())

	s3Service := s3.NewService(appCfg.S3Service, appCfg.S3Bucket)
	routes.RegisterRoutes(r, s3Service, appCfg.S3Bucket)

	log.Printf("Starting server on :%d\n", appCfg.Port)
	if err := r.Run(":" + strconv.Itoa(appCfg.Port)); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
