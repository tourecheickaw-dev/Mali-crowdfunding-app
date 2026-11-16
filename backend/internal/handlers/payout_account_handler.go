package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

type PayoutAccountHandler struct {
	payoutService *service.PayoutAccountService
}

func NewPayoutAccountHandler(payoutService *service.PayoutAccountService) *PayoutAccountHandler {
	return &PayoutAccountHandler{payoutService: payoutService}
}

type addPayoutAccountRequest struct {
	Type          models.PayoutType `json:"type" binding:"required"`
	AccountName   string            `json:"account_name" binding:"required,max=100"`
	AccountNumber string            `json:"account_number" binding:"required,max=50"`
	BankName      string            `json:"bank_name" binding:"max=100"`
	BankCode      string            `json:"bank_code" binding:"max=50"`
	Country       string            `json:"country" binding:"omitempty,len=2"`
}

// Add adds a payout account to the connected organizer.
// To be secured before any deployment: Wave/Orange re-authentication, notification to the old
// number, 24-48 h delay, and 2FA (project security decisions).
func (h *PayoutAccountHandler) Add(c *gin.Context) {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	var req addPayoutAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}

	account, err := h.payoutService.AddPayoutAccount(c.Request.Context(), service.AddPayoutAccountInput{
		OrganizerID:   organizerID,
		Type:          req.Type,
		AccountName:   req.AccountName,
		AccountNumber: req.AccountNumber,
		BankName:      req.BankName,
		BankCode:      req.BankCode,
		Country:       strings.ToUpper(req.Country),
	})
	if err != nil {
		respondServiceError(c, err, "unable to add account")
		return
	}

	c.JSON(http.StatusCreated, account)
}

// List returns the payout accounts of the connected organizer.
func (h *PayoutAccountHandler) List(c *gin.Context) {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	accounts, err := h.payoutService.GetOrganizerPayoutAccounts(c.Request.Context(), organizerID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if accounts == nil {
		accounts = []models.PayoutAccount{}
	}

	c.JSON(http.StatusOK, accounts)
}

// SetActive chooses the account that receives upcoming donations (same protections to plan as Add).
func (h *PayoutAccountHandler) SetActive(c *gin.Context) {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if err := h.payoutService.SetActiveAccount(c.Request.Context(), organizerID, c.Param("id")); err != nil {
		respondServiceError(c, err, "unable to change active account")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "active account updated"})
}
