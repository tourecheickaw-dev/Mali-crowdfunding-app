package handlers

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

func devEnabled() bool {
	return os.Getenv("DEV_AUTH") == "true"
}

type DevHandler struct {
	donationService *service.DonationService
}

func NewDevHandler(donationService *service.DonationService) *DevHandler {
	return &DevHandler{donationService: donationService}
}

// CompleteDonation simulates "the aggregator confirmed the payment", to test the flow locally.
func (h *DevHandler) CompleteDonation(c *gin.Context) {
	if !devEnabled() {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	if err := h.donationService.CompleteDonation(c.Request.Context(), c.Param("id")); err != nil {
		log.Printf("dev complete donation: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()}) // visible detail: dev mode only
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "donation confirmed (simulation)"})
}
