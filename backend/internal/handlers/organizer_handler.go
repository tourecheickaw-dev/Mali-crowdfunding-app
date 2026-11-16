package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

type OrganizerHandler struct {
	organiserService *service.OrganiserService
}

func NewOrganizerHandler(organiserService *service.OrganiserService) *OrganizerHandler {
	return &OrganizerHandler{organiserService: organiserService}
}

type registerOrganizerRequest struct {
	DisplayName string `json:"display_name" binding:"required,max=100"`
	Phone       string `json:"phone" binding:"required,max=30"`
	Email       string `json:"email" binding:"omitempty,email,max=254"`
}

// Register creates an organizer.
// TEMPORARY: without number verification (OTP, Wave, or Orange), anyone could register
// another person's number. Disabled outside dev mode until authentication exists.
func (h *OrganizerHandler) Register(c *gin.Context) {
	if !devEnabled() {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "registration unavailable: authentication not configured"})
		return
	}

	var req registerOrganizerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: name, phone (with country code) and valid email are required"})
		return
	}

	organizer, err := h.organiserService.RegisterOrganizer(c.Request.Context(), service.RegisterOrganizerInput{
		DisplayName: req.DisplayName,
		Phone:       req.Phone,
		Email:       req.Email,
	})
	if err != nil {
		respondServiceError(c, err, "unable to create profile")
		return
	}

	c.JSON(http.StatusCreated, organizer)
}

// publicOrganizer: what is shown to the public. Never phone or email.
type publicOrganizer struct {
	ID                  string `json:"id"`
	DisplayName         string `json:"display_name"`
	Bio                 string `json:"bio,omitempty"`
	ProfilePhotoURL     string `json:"profile_photo_url,omitempty"`
	Country             string `json:"country"`
	Verified            bool   `json:"verified"`
	TotalRaisedLifetime int64  `json:"total_raised_lifetime"`
}

// GetPublic displays the public page of an organizer.
func (h *OrganizerHandler) GetPublic(c *gin.Context) {
	organizer, err := h.organiserService.GetOrganizerByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if organizer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "organizer not found"})
		return
	}

	c.JSON(http.StatusOK, publicOrganizer{
		ID:                  organizer.ID,
		DisplayName:         organizer.DisplayName,
		Bio:                 organizer.Bio,
		ProfilePhotoURL:     organizer.ProfilePhotoURL,
		Country:             organizer.Country,
		Verified:            organizer.VerificationStatus == models.OrganizerVerified,
		TotalRaisedLifetime: organizer.TotalRaisedLifetime,
	})
}
