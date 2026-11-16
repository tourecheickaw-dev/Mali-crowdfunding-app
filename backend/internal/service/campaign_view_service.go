package service

import (
	"context"
	"log"
	"sort"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/repository"
)

// CampaignView is what the public site receives: the campaign plus its cover photo.
type CampaignView struct {
	models.Campaign
	CoverURL       string `json:"cover_url,omitempty"`
	CoverSensitive bool   `json:"cover_sensitive,omitempty"`
}

type CampaignViewService struct {
	campaigns *CampaignService
	media     *repository.MediaRepository
}

func NewCampaignViewService(campaigns *CampaignService, media *repository.MediaRepository) *CampaignViewService {
	return &CampaignViewService{campaigns: campaigns, media: media}
}

// viewOf adds the cover: the image marked as cover, else the first image.
func (s *CampaignViewService) viewOf(ctx context.Context, c models.Campaign) CampaignView {
	view := CampaignView{Campaign: c}

	items, err := s.media.GetMediaByCampaignID(ctx, c.ID)
	if err != nil {
		// a missing photo must never hide the campaign
		log.Printf("load media for campaign %s: %v", c.ID, err)
		return view
	}

	var chosen *models.CampaignMedia
	for i := range items {
		if items[i].MediaType != "image" {
			continue
		}
		if items[i].IsCover {
			chosen = &items[i]
			break
		}
		if chosen == nil {
			chosen = &items[i]
		}
	}
	if chosen != nil {
		view.CoverURL = chosen.URL
		view.CoverSensitive = chosen.IsSensitive
	}
	return view
}

// GetBySlug returns nil when the campaign does not exist or is still a draft.
func (s *CampaignViewService) GetBySlug(ctx context.Context, slug string) (*CampaignView, error) {
	c, err := s.campaigns.GetCampaignBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if c == nil || c.Status == models.CampaignDraft {
		return nil, nil
	}
	v := s.viewOf(ctx, *c)
	return &v, nil
}

// ListActive returns the newest active campaigns, at most limit.
func (s *CampaignViewService) ListActive(ctx context.Context, limit int) ([]CampaignView, error) {
	list, err := s.campaigns.ListActive(ctx)
	if err != nil {
		return nil, err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].CreatedAt.After(list[j].CreatedAt) })
	if len(list) > limit {
		list = list[:limit]
	}
	views := make([]CampaignView, 0, len(list))
	for _, c := range list {
		views = append(views, s.viewOf(ctx, c))
	}
	return views, nil
}
