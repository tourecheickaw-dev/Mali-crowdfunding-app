package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type replyRequest struct {
	Reply string `json:"reply" binding:"required,max=500"`
}

// Reply: the organizer replies to a donor's comment on THEIR campaign.
func (h *CommentHandler) Reply(c *gin.Context) {
	campaign := loadOwnedCampaign(c, h.campaignService)
	if campaign == nil {
		return
	}

	var req replyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid reply (500 characters maximum)"})
		return
	}

	if err := h.donationService.ReplyToComment(c.Request.Context(), campaign.ID, c.Param("donationID"), req.Reply); err != nil {
		respondServiceError(c, err, "failed to save reply")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reply saved"})
}

// DeleteReply: the organizer removes their reply.
func (h *CommentHandler) DeleteReply(c *gin.Context) {
	campaign := loadOwnedCampaign(c, h.campaignService)
	if campaign == nil {
		return
	}

	if err := h.donationService.DeleteReply(c.Request.Context(), campaign.ID, c.Param("donationID")); err != nil {
		respondServiceError(c, err, "failed to remove reply")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "reply removed"})
}
