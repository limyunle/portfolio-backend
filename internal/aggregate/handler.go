package aggregate

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) Refresh(c *gin.Context) {
	username := c.DefaultQuery("username", "limyunle")

	if err := h.Service.RefreshAndStore(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "aggregate stats refreshed"})
}

func (h *Handler) ServeJSON(c *gin.Context) {
	stats, err := h.Service.GetFromS3("aggregate-stats.json")
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, stats)
}
