package github

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(s Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) GetRepos(c *gin.Context) {
	username := c.DefaultQuery("username", "limyunle")
	repos, err := h.service.GetRepos(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Println("Successful call to github")
	log.Printf("repos retrieved: %v", repos)
	c.JSON(http.StatusOK, repos)
}
