package routes

import (
	"github.com/limyunle/portfolio-backend/internal/github"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	githubService := github.NewService()
	githubHandler := github.NewHandler(githubService)

	githubRoutes := r.Group("/github")
	{
		githubRoutes.GET("/repos", githubHandler.GetRepos)
	}

}
