package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

const maxPublicDonations = 100

type DonationHistoryHandler struct {
	campaignService *service.CampaignService
	donationService *service.DonationService
}

func NewDonationHistoryHandler(campaignService *service.CampaignService, donationService *service.DonationService) *DonationHistoryHandler {
	return &DonationHistoryHandler{campaignService: campaignService, donationService: donationService}
}

// publicDonation: what is shown to the public. Never show phone, email, or payment account.
// The comment must be displayed as TEXT by the website, never as HTML.
type publicDonation struct {
	ID            string     `json:"id"`
	DonorName     string     `json:"donor_name"`
	DonorPhotoURL string     `json:"donor_photo_url,omitempty"`
	Comment       string     `json:"comment,omitempty"`
	Amount        int64      `json:"amount"`
	Currency      string     `json:"currency"`
	CreatedAt     time.Time  `json:"created_at"`
	Reply         string     `json:"reply,omitempty"`
	ReplyAt       *time.Time `json:"reply_at,omitempty"`
}

// ListPublic lists the latest confirmed donations for a campaign ("Anonymous" if the donor requested it).
func (h *DonationHistoryHandler) ListPublic(c *gin.Context) {
	campaign, err := h.campaignService.GetCampaignBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if campaign == nil || campaign.Status == models.CampaignDraft {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	donations, err := h.donationService.GetCompletedDonations(c.Request.Context(), campaign.ID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}

	result := make([]publicDonation, 0, len(donations))
	for _, d := range donations {
		name := d.DonorName
		photo := d.DonorPhotoURL
		comment := d.Comment

		if d.IsAnonymous || name == "" {
			name = "Anonymous"
		}
		if d.IsAnonymous {
			photo = "" // never include a photo for an anonymous donation
		}
		if d.CommentHidden {
			comment = "" // hidden by the organizer
		}

		// A reply is only displayed alongside a visible comment.
		reply := ""
		var replyAt *time.Time
		if comment != "" {
			reply = d.Reply
			replyAt = d.ReplyAt
		}

		result = append(result, publicDonation{
			ID:            d.ID,
			DonorName:     name,
			DonorPhotoURL: photo,
			Comment:       comment,
			Reply:         reply,
			ReplyAt:       replyAt,
			Amount:        d.Amount,
			Currency:      d.Currency,
			CreatedAt:     d.CreatedAt,
		})
		if len(result) >= maxPublicDonations {
			break
		}
	}
}
