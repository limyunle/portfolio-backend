package s3

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service    *Service
	BucketName string
}

func NewHandler(service *Service, bucket string) *Handler {
	return &Handler{
		Service:    service,
		BucketName: bucket,
	}
}

func (h *Handler) GetJSON(c *gin.Context) {
	key := c.Param("key")
	var data interface{}

	if err := h.Service.GetJSON(c.Request.Context(), h.BucketName, key, &data); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *Handler) UploadJSON(c *gin.Context) {
	key := c.Param("key")
	var payload interface{}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.Service.UploadJSON(c.Request.Context(), h.BucketName, key, payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "uploaded successfully"})
}
