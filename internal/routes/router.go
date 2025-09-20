package routes

import (
	"github.com/limyunle/portfolio-backend/internal/github"
	"github.com/limyunle/portfolio-backend/internal/leetcode"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	githubService := github.NewService()
	githubHandler := github.NewHandler(githubService)

	leetcodeService := leetcode.NewService()
	leetcodeHandler := leetcode.NewHandler(leetcodeService)

	githubRoutes := r.Group("/github")
	{
		githubRoutes.GET("/repos", githubHandler.GetRepos)
	}

	leetcodeRoutes := r.Group("/leetcode")
	{
		leetcodeRoutes.GET("/stats", leetcodeHandler.GetStats)
	}

}
