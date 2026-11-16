package handlers

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

type CampaignHandler struct {
	campaignService *service.CampaignService
	views           *service.CampaignViewService
}

func NewCampaignHandler(campaignService *service.CampaignService) *CampaignHandler {
	return &CampaignHandler{campaignService: campaignService}
}

func (h *CampaignHandler) SetViewService(v *service.CampaignViewService) { h.views = v }

// createCampaignRequest: only the fields that the organizer is allowed to choose.
// ID, collected amount, status, dates, and slug are determined by the server.
type createCampaignRequest struct {
	Title           string                 `json:"title" binding:"required,max=120"`
	Description     string                 `json:"description" binding:"max=5000"`
	GoalAmount      *int64                 `json:"goal_amount"` // optional
	Country         string                 `json:"country" binding:"required,len=2"`
	BeneficiaryType models.BeneficiaryType `json:"beneficiary_type"`
	BeneficiaryName string                 `json:"beneficiary_name" binding:"max=120"`
}

// Create a campaign for the currently authenticated organizer.
func (h *CampaignHandler) Create(c *gin.Context) {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req createCampaignRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: check the title, country (2 letters), and amounts"})
		return
	}

	if req.GoalAmount != nil && *req.GoalAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the goal amount, if provided, must be greater than zero"})
		return
	}

	switch req.BeneficiaryType {
	case "":
		req.BeneficiaryType = models.BeneficiarySelf
	case models.BeneficiarySelf, models.BeneficiaryOtherPerson, models.BeneficiaryCauseCollective:
		// valid value
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid beneficiary type"})
		return
	}

	campaign, err := h.campaignService.CreateCampaign(c.Request.Context(), service.CreateCampaignInput{
		OrganizerID:     organizerID,
		Title:           req.Title,
		Description:     req.Description,
		GoalAmount:      req.GoalAmount,
		Country:         strings.ToUpper(req.Country),
		BeneficiaryType: req.BeneficiaryType,
		BeneficiaryName: req.BeneficiaryName,
	})
	if err != nil {
		// Keep internal details in logs: the client must not see internal errors (AWS, etc.)
		log.Printf("create campaign: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to create campaign"})
		return
	}

	c.JSON(http.StatusCreated, campaign)
}

// GetBySlug publicly displays a campaign (with its cover photo) using its slug link.
func (h *CampaignHandler) GetBySlug(c *gin.Context) {
	view, err := h.views.GetBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		log.Printf("get campaign by slug: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	if view == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}
	c.JSON(http.StatusOK, view)
}

// ListPublic lists the newest active campaigns for the home and explore pages.
func (h *CampaignHandler) ListPublic(c *gin.Context) {
	list, err := h.views.ListActive(c.Request.Context(), 24)
	if err != nil {
		log.Printf("list public campaigns: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"campaigns": list})
}

// currentOrganizerID returns the connected organizer.
// TEMPORARY, for local testing only: it reads a header that anyone can forge,
// so it is disabled unless DEV_AUTH=true. Must be replaced by real authentication (Cognito) before deployment.
func currentOrganizerID(c *gin.Context) string {
	if os.Getenv("DEV_AUTH") != "true" {
		return ""
	}
	return c.GetHeader("X-Dev-Organizer-ID")
}

// ListMine returns the campaigns of the connected organizer.
func (h *CampaignHandler) ListMine(c *gin.Context) {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	campaigns, err := h.campaignService.GetCampaignsByOrganizer(c.Request.Context(), organizerID)
	if err != nil {
		log.Printf("list my campaigns: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable to load campaigns"})
		return
	}
	if campaigns == nil {
		campaigns = []models.Campaign{} // return [] and not null
	}

	c.JSON(http.StatusOK, gin.H{"campaigns": campaigns})
}
