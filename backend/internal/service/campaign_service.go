package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/repository"
)

type CampaignService struct {
	campaignRepo *repository.CampaignRepository
}

func NewCampaignService(campaignRepo *repository.CampaignRepository) *CampaignService {
	return &CampaignService{campaignRepo: campaignRepo}
}

type CreateCampaignInput struct {
	OrganizerID     string
	Title           string
	Description     string
	GoalAmount      *int64 // nil = no target set, product decision already finalized
	Currency        string
	Country         string
	BeneficiaryType models.BeneficiaryType
	BeneficiaryName string
}

func (s *CampaignService) CreateCampaign(ctx context.Context, input CreateCampaignInput) (*models.Campaign, error) {
	if input.OrganizerID == "" {
		return nil, errors.New("organizer ID is required")
	}
	if input.Title == "" {
		return nil, errors.New("campaign title is required")
	}
	if input.GoalAmount != nil && *input.GoalAmount <= 0 {
		return nil, errors.New("goal amount, if provided, must be greater than zero")
	}

	if input.Currency == "" {
		input.Currency = "XOF"
	}

	campaignID := uuid.New().String()
	slug := generateSlug(input.Title)

	existing, err := s.campaignRepo.GetCampaignBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("failed to check slug uniqueness: %w", err)
	}
	if existing != nil {
		slug = fmt.Sprintf("%s-%s", slug, campaignID[:8])
	}

	now := time.Now()
	campaign := models.Campaign{
		ID:              campaignID,
		OrganizerID:     input.OrganizerID,
		Title:           input.Title,
		Description:     input.Description,
		GoalAmount:      input.GoalAmount,
		CurrentAmount:   0,
		Currency:        input.Currency,
		Country:         input.Country,
		BeneficiaryType: input.BeneficiaryType,
		BeneficiaryName: input.BeneficiaryName,
		URLSlug:         slug,
		Status:          models.CampaignActive,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.campaignRepo.CreateCampaign(ctx, campaign); err != nil {
		return nil, fmt.Errorf("failed to create campaign: %w", err)
	}

	return &campaign, nil
}

func (s *CampaignService) GetCampaignByID(ctx context.Context, id string) (*models.Campaign, error) {
	if id == "" {
		return nil, errors.New("campaign ID is required")
	}
	return s.campaignRepo.GetCampaignByID(ctx, id)
}

func (s *CampaignService) GetCampaignBySlug(ctx context.Context, slug string) (*models.Campaign, error) {
	if slug == "" {
		return nil, errors.New("campaign slug is required")
	}
	return s.campaignRepo.GetCampaignBySlug(ctx, slug)
}

func (s *CampaignService) GetCampaignsByOrganizer(ctx context.Context, organizerID string) ([]models.Campaign, error) {
	if organizerID == "" {
		return nil, errors.New("organizer ID is required")
	}
	return s.campaignRepo.GetCampaignsByOrganizer(ctx, organizerID)
}

var slugReplacer = strings.NewReplacer(
	"à", "a", "â", "a", "ä", "a", "ç", "c",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"î", "i", "ï", "i", "ô", "o", "ö", "o",
	"ù", "u", "û", "u", "ü", "u", "ÿ", "y",
	"œ", "oe", "æ", "ae",
)

func generateSlug(title string) string {
	s := slugReplacer.Replace(strings.ToLower(title))
	var result []rune
	for _, r := range s {
		switch {
		case (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'):
			result = append(result, r)
		case r == ' ' || r == '-' || r == '_' || r == '\'' || r == '’':
			if len(result) > 0 && result[len(result)-1] != '-' {
				result = append(result, '-')
			}
		}
	}
	if len(result) > 60 {
		result = result[:60]
	}
	slug := strings.Trim(string(result), "-")
	if slug == "" {
		slug = "cagnotte" // title with no latin letters
	}
	return slug

}

func (s *CampaignService) ListActive(ctx context.Context) ([]models.Campaign, error) {
	return s.campaignRepo.ListActiveCampaigns(ctx, 200)
}
