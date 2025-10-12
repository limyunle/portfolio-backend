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

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	rg.GET("/refresh", h.Refresh)
}

func (h *Handler) Refresh(c *gin.Context) {
	username := c.DefaultQuery("username", "limyunle")

	if err := h.Service.RefreshAndStore(username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "aggregate stats refreshed"})
}
