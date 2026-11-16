package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

// GetCompletedDonations returns the completed donations of a campaign, from most recent to oldest.
func (s *DonationService) GetCompletedDonations(ctx context.Context, campaignID string) ([]models.Donation, error) {
	if campaignID == "" {
		return nil, errors.New("campaign ID is required")
	}

	all, err := s.donationRepo.GetDonationsByCampaignID(ctx, campaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donations: %w", err)
	}

	completed := make([]models.Donation, 0, len(all))
	for _, d := range all {
		if d.Status == models.DonationCompleted {
			completed = append(completed, d)
		}
	}

	sort.Slice(completed, func(i, j int) bool {
		return completed[i].CreatedAt.After(completed[j].CreatedAt)
	})

	return completed, nil
}

// GetDonationsForDonor returns a donor's donations, from newest to oldest.
func (s *DonationService) GetDonationsForDonor(ctx context.Context, donorID string) ([]models.Donation, error) {
	if donorID == "" {
		return nil, errors.New("donor ID is required")
	}
	donations, err := s.donationRepo.GetDonationsByDonorID(ctx, donorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donor donations: %w", err)
	}
	sort.Slice(donations, func(i, j int) bool {
		return donations[i].CreatedAt.After(donations[j].CreatedAt)
	})
	return donations, nil
}
