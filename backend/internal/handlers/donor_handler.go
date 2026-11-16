package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

const maxDonorDonations = 200

type DonorHandler struct {
	donorService    *service.DonorService
	donationService *service.DonationService
}

func NewDonorHandler(donorService *service.DonorService, donationService *service.DonationService) *DonorHandler {
	return &DonorHandler{donorService: donorService, donationService: donationService}
}

// currentDonorID returns the connected donor.
// TEMPORARY, for local testing only: the header is forgeable, so disabled
// unless dev mode is enabled. To be replaced with actual donor verification (OTP, Wave, or Orange).
func currentDonorID(c *gin.Context) string {
	if !devEnabled() {
		return ""
	}
	return c.GetHeader("X-Dev-Donor-ID")
}

// donorProfile: the donor sees THEIR own data, never someone else's.
type donorProfile struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Phone           string `json:"phone,omitempty"`
	Email           string `json:"email,omitempty"`
	Verified        bool   `json:"verified"`
	ProfilePhotoURL string `json:"profile_photo_url,omitempty"`
}

// GetProfile returns the profile of the connected donor.
func (h *DonorHandler) GetProfile(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	donor, err := h.donorService.GetDonorByID(c.Request.Context(), donorID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if donor == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "profile not found"})
		return
	}

	c.JSON(http.StatusOK, donorProfile{
		ID:              donor.ID,
		Name:            donor.Name,
		Phone:           donor.Phone,
		Email:           donor.Email,
		ProfilePhotoURL: donor.ProfilePhotoURL,
		Verified:        donor.VerificationStatus == models.DonorVerified,
	})
}

// donorDonation: a donation seen by its donor. Neither organizer payout account
// nor internal transfer reference.
type donorDonation struct {
	ID            string                `json:"id"`
	CampaignID    string                `json:"campaign_id"`
	Amount        int64                 `json:"amount"`
	FeeAmount     int64                 `json:"fee_amount"`
	TipAmount     int64                 `json:"tip_amount"`
	TotalCharged  int64                 `json:"total_charged"`
	Currency      string                `json:"currency"`
	PaymentMethod models.PaymentMethod  `json:"payment_method"`
	Status        models.DonationStatus `json:"status"`
	Comment       string                `json:"comment,omitempty"`
	CommentHidden bool                  `json:"comment_hidden,omitempty"`
	IsAnonymous   bool                  `json:"is_anonymous"`
	CreatedAt     time.Time             `json:"created_at"`
	Reply         string                `json:"reply,omitempty"`
}

// ListDonations returns the donations of the connected donor, from newest to oldest.
func (h *DonorHandler) ListDonations(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	donations, err := h.donationService.GetDonationsForDonor(c.Request.Context(), donorID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}

	result := make([]donorDonation, 0, len(donations))
	for _, d := range donations {
		result = append(result, donorDonation{
			ID:            d.ID,
			CampaignID:    d.CampaignID,
			Amount:        d.Amount,
			FeeAmount:     d.FeeAmount,
			TipAmount:     d.TipAmount,
			TotalCharged:  d.TotalCharged,
			Currency:      d.Currency,
			PaymentMethod: d.PaymentMethod,
			Status:        d.Status,
			IsAnonymous:   d.IsAnonymous,
			Comment:       d.Comment,
			CommentHidden: d.CommentHidden,
			CreatedAt:     d.CreatedAt,
			Reply:         d.Reply,
		})
		if len(result) >= maxDonorDonations {
			break
		}
	}

	c.JSON(http.StatusOK, result)
}

type setPhotoRequest struct {
	URL string `json:"url" binding:"required,max=2048"`
}

// SetPhoto sets the profile photo of the logged-in donor.
func (h *DonorHandler) SetPhoto(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req setPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid photo link"})
		return
	}

	h.applyPhoto(c, donorID, req.URL)
}

// RemovePhoto removes the profile photo of the logged-in donor.
func (h *DonorHandler) RemovePhoto(c *gin.Context) {
	donorID := currentDonorID(c)
	if donorID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	h.applyPhoto(c, donorID, "")
}

// applyPhoto saves the photo, then aligns the public history of their donations.
func (h *DonorHandler) applyPhoto(c *gin.Context, donorID, url string) {
	donor, err := h.donorService.SetProfilePhoto(c.Request.Context(), donorID, url)
	if err != nil {
		respondServiceError(c, err, "unable to update photo")
		return
	}

	// Photo displayed publicly on their donations: only if the donor is verified.
	visible := ""
	if donor.VerificationStatus == models.DonorVerified {
		visible = donor.ProfilePhotoURL
	}
	if err := h.donationService.SyncDonorPhoto(c.Request.Context(), donorID, visible); err != nil {
		respondServiceError(c, err, "photo saved, but history update failed")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "photo updated"})
}
