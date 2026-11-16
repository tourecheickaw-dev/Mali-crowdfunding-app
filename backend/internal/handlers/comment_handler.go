package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

const maxOrganizerDonations = 200

type CommentHandler struct {
	campaignService *service.CampaignService
	donationService *service.DonationService
}

func NewCommentHandler(campaignService *service.CampaignService, donationService *service.DonationService) *CommentHandler {
	return &CommentHandler{campaignService: campaignService, donationService: donationService}
}

// loadOwnedCampaign returns the campaign from the URL (:id) if it belongs to the logged-in organizer.
// Otherwise responds directly (401 or 404) and returns nil.
func loadOwnedCampaign(c *gin.Context, campaignService *service.CampaignService) *models.Campaign {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return nil
	}

	campaign, err := campaignService.GetCampaignByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		respondServiceError(c, err, "server error")
		return nil
	}
	if campaign == nil || campaign.OrganizerID != organizerID {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return nil
	}
	return campaign
}

type commentRequest struct {
	Comment string `json:"comment" binding:"required,max=500"`
}

// UpdateOwn: the logged-in donor updates their comment.
func (h *CommentHandler) UpdateOwn(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req commentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid comment (maximum 500 characters)"})
		return
	}

	if err := h.donationService.UpdateComment(c.Request.Context(), donorID, c.Param("id"), req.Comment); err != nil {
		respondServiceError(c, err, "unable to update comment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment updated"})
}

// DeleteOwn: the logged-in donor deletes their comment.
func (h *CommentHandler) DeleteOwn(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if err := h.donationService.DeleteComment(c.Request.Context(), donorID, c.Param("id")); err != nil {
		respondServiceError(c, err, "unable to delete comment")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment deleted"})
}

// organizerDonationView: a donation viewed by the campaign organizer (to moderate comments).
// Neither phone nor email.
type organizerDonationView struct {
	ID            string     `json:"id"`
	DonorName     string     `json:"donor_name"`
	IsAnonymous   bool       `json:"is_anonymous"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	Comment       string     `json:"comment,omitempty"`
	CommentHidden bool       `json:"comment_hidden"`
	CreatedAt     time.Time  `json:"created_at"`
	Reply         string     `json:"reply,omitempty"`
	ReplyAt       *time.Time `json:"reply_at,omitempty"`
}

// ListForOrganizer: the confirmed donations of THEIR campaign, including hidden comments.
func (h *CommentHandler) ListForOrganizer(c *gin.Context) {
	campaign := loadOwnedCampaign(c, h.campaignService)
	if campaign == nil {
		return
	}

	donations, err := h.donationService.GetCompletedDonations(c.Request.Context(), campaign.ID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}

	result := make([]organizerDonationView, 0, len(donations))
	for _, d := range donations {
		result = append(result, organizerDonationView{
			ID:            d.ID,
			DonorName:     d.DonorName,
			IsAnonymous:   d.IsAnonymous,
			Amount:        d.Amount,
			Currency:      d.Currency,
			Comment:       d.Comment,
			CommentHidden: d.CommentHidden,
			CreatedAt:     d.CreatedAt,
			Reply:         d.Reply,
			ReplyAt:       d.ReplyAt,
		})
		if len(result) >= maxOrganizerDonations {
			break
		}
	}

	c.JSON(http.StatusOK, result)
}

type visibilityRequest struct {
	Hidden *bool `json:"hidden" binding:"required"`
}

// SetVisibility: the organizer hides or unhides a comment on THEIR campaign.
func (h *CommentHandler) SetVisibility(c *gin.Context) {
	campaign := loadOwnedCampaign(c, h.campaignService)
	if campaign == nil {
		return
	}

	var req visibilityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "indicate \"hidden\": true or false"})
		return
	}

	if err := h.donationService.SetCommentHidden(c.Request.Context(), campaign.ID, c.Param("donationID"), *req.Hidden); err != nil {
		respondServiceError(c, err, "unable to change comment visibility")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "comment visibility updated"})
}
