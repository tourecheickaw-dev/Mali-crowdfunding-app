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

type OrganiserService struct {
	organiserRepo *repository.OrganizerRepository
}

func NewOrganiserService(organiserRepo *repository.OrganizerRepository) *OrganiserService {
	return &OrganiserService{
		organiserRepo: organiserRepo,
	}
}

// RegisterOrganizerInput holds data required to register a new campaign organizer
type RegisterOrganizerInput struct {
	DisplayName string `json:"display_name"`
	Phone       string `json:"phone"`
	Email       string `json:"email,omitempty"`
}

// RegisterOrganizer normalizes the phone, validates the UEMOA zone, and creates an organizer record.
func (s *OrganiserService) RegisterOrganizer(ctx context.Context, input RegisterOrganizerInput) (*models.Organizer, error) {
	if input.DisplayName == "" {
		return nil, errors.New("display name is required")
	}

	phone, err := normalizePhone(input.Phone)
	if err != nil {
		return nil, fmt.Errorf("invalid phone number: %w", err)
	}

	// organizer phone number must belong to the UEMOA zone
	if !isUEMOAPhone(phone) {
		return nil, errors.New("organizer phone number must belong to the UEMOA zone (e.g., Mali +223, Senegal +221, etc.)")
	}

	existing, err := s.organiserRepo.GetOrganizerByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing organizer: %w", err)
	}
	if existing != nil {
		return nil, errors.New("an organizer with this phone number already exists")
	}

	now := time.Now()
	newOrganizer := models.Organizer{
		ID:                 uuid.New().String(),
		DisplayName:        input.DisplayName,
		Phone:              phone,
		Email:              input.Email,
		Country:            countryFromPhone(phone),
		VerificationStatus: models.OrganizerUnverified,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.organiserRepo.SaveOrganizer(ctx, &newOrganizer); err != nil {
		return nil, fmt.Errorf("failed to save organizer: %w", err)
	}

	return &newOrganizer, nil
}

// GetOrganizerByID retrieves an organizer profile by ID
func (s *OrganiserService) GetOrganizerByID(ctx context.Context, id string) (*models.Organizer, error) {
	if id == "" {
		return nil, errors.New("organizer ID is required")
	}
	return s.organiserRepo.GetOrganizerByID(ctx, id)
}

// GetOrganizerByPhone retrieves an organizer profile by normalized phone number
func (s *OrganiserService) GetOrganizerByPhone(ctx context.Context, phone string) (*models.Organizer, error) {
	normalized, err := normalizePhone(phone)
	if err != nil {
		return nil, err
	}
	return s.organiserRepo.GetOrganizerByPhone(ctx, normalized)
}
