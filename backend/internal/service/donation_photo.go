package service

import (
	"context"
	"errors"
	"fmt"
)

// SyncDonorPhoto aligns the photo displayed on all past donations of a donor
// (empty photoURL = no more photo). Called when they change or remove their photo,
// so that a removed photo also disappears from the public history.
func (s *DonationService) SyncDonorPhoto(ctx context.Context, donorID, photoURL string) error {
	if donorID == "" {
		return errors.New("donor ID is required")
	}

	donations, err := s.donationRepo.GetDonationsByDonorID(ctx, donorID)
	if err != nil {
		return fmt.Errorf("failed to fetch donor donations: %w", err)
	}

	for _, d := range donations {
		wanted := photoURL
		if d.IsAnonymous {
			wanted = "" // never show photo on an anonymous donation
		}
		if d.DonorPhotoURL == wanted {
			continue
		}
		if err := s.donationRepo.UpdateDonorPhoto(ctx, d.ID, wanted); err != nil {
			return fmt.Errorf("failed to update photo of donation %s: %w", d.ID, err)
		}
	}

	return nil
}
