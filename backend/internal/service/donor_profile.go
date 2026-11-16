package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

// SetProfilePhoto sets the donor's profile photo (empty url = remove it).
// For now, an https link is accepted; eventually, an upload to our own storage (S3).
func (s *DonorService) SetProfilePhoto(ctx context.Context, donorID, url string) (*models.Donor, error) {
	if donorID == "" {
		return nil, errors.New("donor ID is required")
	}
	if url != "" && !strings.HasPrefix(url, "https://") {
		return nil, errors.New("photo link must start with https://")
	}
	if len(url) > 2048 {
		return nil, errors.New("photo link is too long")
	}

	donor, err := s.donorRepo.GetDonorByID(ctx, donorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donor: %w", err)
	}
	if donor == nil {
		return nil, errors.New("donor not found")
	}

	donor.ProfilePhotoURL = url
	donor.UpdatedAt = time.Now()
	if err := s.donorRepo.SaveDonor(ctx, donor); err != nil {
		return nil, fmt.Errorf("failed to save donor: %w", err)
	}

	return donor, nil
}
