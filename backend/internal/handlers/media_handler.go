package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/service"
)

type MediaHandler struct {
	mediaService    *service.MediaService
	campaignService *service.CampaignService
	uploads         *service.UploadService
}

func NewMediaHandler(mediaService *service.MediaService, campaignService *service.CampaignService) *MediaHandler {
	return &MediaHandler{mediaService: mediaService, campaignService: campaignService}
	// SetUploadService enables the "our bucket only" rule for media URLs.
	// SetUploadService enables the "our bucket only" rule for media URLs.

}

// SetUploadService enables the "our bucket only" rule for media URLs.
func (h *MediaHandler) SetUploadService(u *service.UploadService) {
	h.uploads = u
}

type attachMediaRequest struct {
	URL          string           `json:"url" binding:"required,max=2048"`
	MediaType    models.MediaType `json:"media_type" binding:"required"`
	Caption      string           `json:"caption" binding:"max=300"`
	IsCover      bool             `json:"is_cover"`
	IsSensitive  bool             `json:"is_sensitive"`
	DisplayOrder int              `json:"display_order" binding:"gte=0,lte=1000"`
}

// ownedCampaign returns the campaign from the URL if it belongs to the connected organizer.
// Otherwise, it replies directly (401 or 404) and returns nil.
func (h *MediaHandler) ownedCampaign(c *gin.Context) *models.Campaign {
	organizerID := currentOrganizerID(c)
	if organizerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return nil
	}

	campaign, err := h.campaignService.GetCampaignByID(c.Request.Context(), c.Param("id"))
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

// Attach adds a photo or a video to a campaign.
// Eventually, files will be sent to our storage (S3) and only those links will be accepted.
// In the meantime, https is required at a minimum.
func (h *MediaHandler) Attach(c *gin.Context) {
	campaign := h.ownedCampaign(c)
	if campaign == nil {
		return
	}

	var req attachMediaRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		return
	}
	if h.uploads != nil {
		if !h.uploads.IsOwnedURL(campaign.OrganizerID, req.URL) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "the photo must be uploaded through the app"})
			return
		}
	} else if !strings.HasPrefix(req.URL, "https://") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "the link must start with https://"})
		return

	}

	media, err := h.mediaService.AttachMediaToCampaign(c.Request.Context(), service.AttachMediaInput{
		CampaignID:   campaign.ID,
		URL:          req.URL,
		MediaType:    req.MediaType,
		Caption:      req.Caption,
		IsCover:      req.IsCover,
		IsSensitive:  req.IsSensitive,
		DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		respondServiceError(c, err, "unable to add media")
		return
	}

	c.JSON(http.StatusCreated, media)
}

// ListPublic returns the media of a public campaign. is_sensitive allows the site to blur them.
func (h *MediaHandler) ListPublic(c *gin.Context) {
	campaign, err := h.campaignService.GetCampaignBySlug(c.Request.Context(), c.Param("slug"))
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if campaign == nil || campaign.Status == models.CampaignDraft {
		c.JSON(http.StatusNotFound, gin.H{"error": "campaign not found"})
		return
	}

	items, err := h.mediaService.GetMediaByCampaignID(c.Request.Context(), campaign.ID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}
	if items == nil {
		items = []models.CampaignMedia{}
	}

	c.JSON(http.StatusOK, items)
}

// SetCover sets an existing media as cover.
func (h *MediaHandler) SetCover(c *gin.Context) {
	campaign := h.ownedCampaign(c)
	if campaign == nil {
		return
	}

	if err := h.mediaService.SetCoverMedia(c.Request.Context(), campaign.ID, c.Param("mediaID")); err != nil {
		respondServiceError(c, err, "unable to change cover")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "cover updated"})
}

// Delete deletes a media, after verifying that it belongs to this campaign.
func (h *MediaHandler) Delete(c *gin.Context) {
	campaign := h.ownedCampaign(c)
	if campaign == nil {
		return
	}

	items, err := h.mediaService.GetMediaByCampaignID(c.Request.Context(), campaign.ID)
	if err != nil {
		respondServiceError(c, err, "server error")
		return
	}

	mediaID := c.Param("mediaID")
	found := false
	for _, m := range items {
		if m.ID == mediaID {
			found = true
			break
		}
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "media not found"})
		return
	}

	if err := h.mediaService.DeleteMedia(c.Request.Context(), mediaID); err != nil {
		respondServiceError(c, err, "unable to delete media")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "media deleted"})
}
