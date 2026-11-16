package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func respondServiceError(c *gin.Context, err error, defaultMsg string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": defaultMsg})
}
