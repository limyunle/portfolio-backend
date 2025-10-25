package routes

import (
	"github.com/limyunle/portfolio-backend/internal/aggregate"
	"github.com/limyunle/portfolio-backend/internal/github"
	"github.com/limyunle/portfolio-backend/internal/leetcode"
	"github.com/limyunle/portfolio-backend/internal/s3"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, s3Client *s3.Service, bucketName string) {
	// --------------------
	// GitHub Routes, for internal testing only
	// --------------------
	githubService := github.NewService()
	githubHandler := github.NewHandler(githubService)

	githubRoutes := r.Group("/github")
	{
		githubRoutes.GET("/repos", githubHandler.GetRepos)
	}

	// --------------------
	// LeetCode Routes, for internal testing only
	// --------------------
	leetcodeService := leetcode.NewService()
	leetcodeHandler := leetcode.NewHandler(leetcodeService)

	leetcodeRoutes := r.Group("/leetcode")
	{
		leetcodeRoutes.GET("/stats", leetcodeHandler.GetLeetCodeStats)
	}

	// --------------------
	// S3 Routes
	// --------------------
	s3Handler := s3.NewHandler(s3Client, bucketName)
	s3Routes := r.Group("/s3")
	{
		s3Routes.GET("/get/:key", s3Handler.GetJSON)
		s3Routes.POST("/upload/:key", s3Handler.UploadJSON)
	}

	// --------------------
	// Aggregate Routes
	// --------------------
	aggregateService := &aggregate.Service{
		GitHubService:   githubService,
		LeetCodeService: leetcodeService,
		S3Service:       s3Client,
		BucketName:      bucketName,
	}
	aggregateHandler := aggregate.NewHandler(aggregateService)

	aggregateRoutes := r.Group("/aggregate")
	{
		aggregateRoutes.GET("/refresh", aggregateHandler.Refresh)
		aggregateRoutes.GET("/stats", aggregateHandler.ServeJSON)
	}
}
