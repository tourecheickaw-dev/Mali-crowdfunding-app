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

type DonorService struct {
	donorRepo *repository.DonorRepository
}

func NewDonorService(donorRepo *repository.DonorRepository) *DonorService {
	return &DonorService{
		donorRepo: donorRepo,
	}
}

// CreateOrGetDonorInput holds data needed to identify or register a donor
type CreateOrGetDonorInput struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Email string `json:"email,omitempty"`
}

// RegisterOrGetDonor checks if a donor exists by phone number; if not, creates a new donor record.
// Donors from all countries are accepted (diaspora), provided they have a country code.
// ATTENTION: The returned donor contains personal data. The handler must never
// return it as-is to the client until the phone number has been verified.
func (s *DonorService) RegisterOrGetDonor(ctx context.Context, input CreateOrGetDonorInput) (*models.Donor, error) {
	phone, err := normalizePhone(input.Phone)
	if err != nil {
		return nil, err
	}

	existingDonor, err := s.donorRepo.GetDonorByPhone(ctx, phone)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing donor by phone: %w", err)
	}
	if existingDonor != nil {
		return existingDonor, nil
	}

	name := input.Name
	if name == "" {
		name = "Anonymous"
	}

	now := time.Now()
	newDonor := models.Donor{
		ID:                 uuid.New().String(),
		Name:               name,
		Phone:              phone,
		Email:              input.Email,
		VerificationMethod: models.DonorVerifNone,
		VerificationStatus: models.DonorUnverified,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.donorRepo.SaveDonor(ctx, &newDonor); err != nil {
		return nil, fmt.Errorf("failed to create donor: %w", err)
	}

	return &newDonor, nil
}

// GetDonorByID retrieves a donor by unique ID
func (s *DonorService) GetDonorByID(ctx context.Context, id string) (*models.Donor, error) {
	if id == "" {
		return nil, errors.New("donor ID is required")
	}
	return s.donorRepo.GetDonorByID(ctx, id)
}

// GetDonorByPhone retrieves a donor by phone number
func (s *DonorService) GetDonorByPhone(ctx context.Context, phone string) (*models.Donor, error) {
	normalized, err := normalizePhone(phone)
	if err != nil {
		return nil, err
	}
	return s.donorRepo.GetDonorByPhone(ctx, normalized)
}
