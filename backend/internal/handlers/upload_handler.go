package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

type UploadHandler struct {
	svc *service.UploadService
}

func NewUploadHandler(svc *service.UploadService) *UploadHandler {
	return &UploadHandler{svc: svc}
}

type presignRequest struct {
	ContentType string `json:"content_type" binding:"required"`
	Size        int64  `json:"size" binding:"required,gt=0"`
}

// Presign prepares a temporary upload URL for one image.
func (h *UploadHandler) Presign(c *gin.Context) {
	ownerID := currentOrganizerID(c)
	if ownerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req presignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	res, err := h.svc.PresignImage(c.Request.Context(), ownerID, req.ContentType, req.Size)
	if err != nil {
		respondServiceError(c, err, "unable to prepare the upload")
		return
	}
	c.JSON(http.StatusOK, res)
}
