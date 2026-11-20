package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

// Known per-transaction limits (in FCFA), to be confirmed with the aggregator. 0 = no known limit.
// This list is also used to reject unknown payment methods.
var maxAmountByMethod = map[models.PaymentMethod]int64{
	models.MethodOrangeMoney: 1000000,
	models.MethodMoov:        500000,
	models.MethodWave:        0,
	models.MethodMTN:         0,
	models.MethodFreeMoney:   0,
	models.MethodWizall:      0,
	models.MethodCard:        0,
}

type DonationHandler struct {
	donationService *service.DonationService
	donorService    *service.DonorService
}

func NewDonationHandler(donationService *service.DonationService, donorService *service.DonorService) *DonationHandler {
	return &DonationHandler{donationService: donationService, donorService: donorService}
}

type createDonationRequest struct {
	CampaignID    string               `json:"campaign_id" binding:"required"`
	Amount        int64                `json:"amount" binding:"required,gt=0,lte=1000000000"` // amount received by the organizer
	TipAmount     int64                `json:"tip_amount" binding:"gte=0,lte=1000000000"`     // optional tip
	PaymentMethod models.PaymentMethod `json:"payment_method" binding:"required"`
	DonorName     string               `json:"donor_name" binding:"max=100"`
	DonorPhone    string               `json:"donor_phone" binding:"omitempty,max=30"`
	DonorEmail    string               `json:"donor_email" binding:"omitempty,email,max=254"`
	Comment       string               `json:"comment" binding:"max=500"`
	IsAnonymous   bool                 `json:"is_anonymous"`
}

// Create records a PENDING donation. It only becomes "completed" upon confirmation by
// the aggregator (signed webhook), never based on the client's word.
// This step does not yet trigger the payment at the aggregator: it will be added with it.
func (h *DonationHandler) Create(c *gin.Context) {
	var req createDonationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data: campaign, amount and payment method required"})
		return
	}

	limit, known := maxAmountByMethod[req.PaymentMethod]
	if !known {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown payment method"})
		return
	}
	// The limit applies to what the donor actually pays: donation + fee + tip.
	totalToPay := req.Amount + service.ComputeFee(req.Amount) + req.TipAmount
	if limit > 0 && totalToPay > limit {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount exceeds the limit for this payment method (fee included)"})
		return
	}

	var donorID *string
	donorPhone := ""
	donorPhoto := ""
	if req.DonorPhone != "" {
		donor, err := h.donorService.RegisterOrGetDonor(c.Request.Context(), service.CreateOrGetDonorInput{
			Name:  req.DonorName,
			Phone: req.DonorPhone,
			Email: req.DonorEmail,
		})
		if err != nil {
			respondServiceError(c, err, "unable to register donor")
			return
		}
		donorID = &donor.ID
		donorPhone = donor.Phone // normalized number

		// Photo displayed only for a verified donor: without verification, anyone
		// could enter another person's number and display their photo next to their comment.
		if donor.VerificationStatus == models.DonorVerified {
			donorPhoto = donor.ProfilePhotoURL
		}
	}

	donation, err := h.donationService.CreateDonation(c.Request.Context(), service.CreateDonationInput{
		CampaignID:    req.CampaignID,
		DonorID:       donorID,
		DonorName:     req.DonorName,
		DonorPhone:    donorPhone,
		DonorEmail:    req.DonorEmail,
		Amount:        req.Amount,
		TipAmount:     req.TipAmount,
		Comment:       req.Comment,
		DonorPhotoURL: donorPhoto,
		PaymentMethod: req.PaymentMethod,
		IsAnonymous:   req.IsAnonymous,
	})
	if err != nil {
		respondServiceError(c, err, "unable to save donation")
		return
	}

	// Deliberately limited response: neither phone, nor email, nor payment account.
	c.JSON(http.StatusCreated, gin.H{
		"id":            donation.ID,
		"campaign_id":   donation.CampaignID,
		"amount":        donation.Amount,
		"fee_amount":    donation.FeeAmount,
		"tip_amount":    donation.TipAmount,
		"total_charged": donation.TotalCharged,
		"currency":      donation.Currency,
		"status":        donation.Status,
	})
}
