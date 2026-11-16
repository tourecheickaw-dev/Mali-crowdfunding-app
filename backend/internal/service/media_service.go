package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/repository"
)

type MediaService struct {
	mediaRepo *repository.MediaRepository
}

func NewMediaService(mediaRepo *repository.MediaRepository) *MediaService {
	return &MediaService{
		mediaRepo: mediaRepo,
	}
}

// AttachMediaInput holds input required to save media metadata
type AttachMediaInput struct {
	CampaignID   string           `json:"campaign_id"`
	URL          string           `json:"url"`
	MediaType    models.MediaType `json:"media_type"`
	Caption      string           `json:"caption,omitempty"`
	IsCover      bool             `json:"is_cover"`
	IsSensitive  bool             `json:"is_sensitive"`
	DisplayOrder int              `json:"display_order,omitempty"`
}

// AttachMediaToCampaign attaches a new media item to a campaign.
// If the media is marked as cover, previous covers will be unset.
func (s *MediaService) AttachMediaToCampaign(ctx context.Context, input AttachMediaInput) (*models.CampaignMedia, error) {
	if input.CampaignID == "" {
		return nil, errors.New("campaign ID is required")
	}
	if input.URL == "" {
		return nil, errors.New("media URL is required")
	}
	if input.MediaType != models.MediaTypeImage && input.MediaType != models.MediaTypeVideo {
		return nil, errors.New("media type must be 'image' or 'video'")
	}

	now := time.Now()
	media := models.CampaignMedia{
		ID:           uuid.New().String(),
		CampaignID:   input.CampaignID,
		URL:          input.URL,
		MediaType:    input.MediaType,
		Caption:      input.Caption,
		IsCover:      input.IsCover,
		IsSensitive:  input.IsSensitive,
		DisplayOrder: input.DisplayOrder,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// We save the new media first: if the remaining steps fail,
	// the campaign still has at least one cover.
	if err := s.mediaRepo.SaveMedia(ctx, &media); err != nil {
		return nil, fmt.Errorf("failed to attach media to campaign: %w", err)
	}

	if media.IsCover {
		if err := s.unsetOtherCovers(ctx, media.CampaignID, media.ID); err != nil {
			return nil, err
		}
	}

	return &media, nil
}

// SetCoverMedia sets an existing media item as the campaign cover.
// It also verifies that the media actually belongs to this campaign.
func (s *MediaService) SetCoverMedia(ctx context.Context, campaignID, mediaID string) error {
	if campaignID == "" {
		return errors.New("campaign ID is required")
	}
	if mediaID == "" {
		return errors.New("media ID is required")
	}

	existing, err := s.mediaRepo.GetMediaByCampaignID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to fetch campaign media: %w", err)
	}

	var target *models.CampaignMedia
	for i := range existing {
		if existing[i].ID == mediaID {
			target = &existing[i]
			break
		}
	}
	if target == nil {
		return errors.New("media not found in this campaign")
	}

	target.IsCover = true
	target.UpdatedAt = time.Now()
	if err := s.mediaRepo.SaveMedia(ctx, target); err != nil {
		return fmt.Errorf("failed to set cover media: %w", err)
	}

	return s.unsetOtherCovers(ctx, campaignID, mediaID)
}

// unsetOtherCovers removes the "cover" status from all media in the campaign,
// except for the one whose ID matches the keepID parameter.
func (s *MediaService) unsetOtherCovers(ctx context.Context, campaignID, keepID string) error {
	existing, err := s.mediaRepo.GetMediaByCampaignID(ctx, campaignID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing media: %w", err)
	}

	for i := range existing {
		m := existing[i]
		if m.IsCover && m.ID != keepID {
			m.IsCover = false
			m.UpdatedAt = time.Now()
			if err := s.mediaRepo.SaveMedia(ctx, &m); err != nil {
				return fmt.Errorf("failed to unset previous cover %s: %w", m.ID, err)
			}
		}
	}

	return nil
}

// GetMediaByCampaignID retrieves all media records associated with a specific campaign
func (s *MediaService) GetMediaByCampaignID(ctx context.Context, campaignID string) ([]models.CampaignMedia, error) {
	if campaignID == "" {
		return nil, errors.New("campaign ID is required")
	}
	return s.mediaRepo.GetMediaByCampaignID(ctx, campaignID)
}

// DeleteMedia removes a media item by ID
func (s *MediaService) DeleteMedia(ctx context.Context, id string) error {
	if id == "" {
		return errors.New("media ID is required")
	}
	return s.mediaRepo.DeleteMedia(ctx, id)
}
