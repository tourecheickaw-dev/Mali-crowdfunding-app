package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

// donationOwnedBy returns a donation only if it belongs to this donor.
// Same response whether the donation does not exist or belongs to someone else.
func (s *DonationService) donationOwnedBy(ctx context.Context, donorID, donationID string) (*models.Donation, error) {
	if donorID == "" || donationID == "" {
		return nil, errors.New("donor ID and donation ID are required")
	}

	donation, err := s.donationRepo.GetDonationByID(ctx, donationID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donation: %w", err)
	}
	if donation == nil || donation.DonorID == nil || *donation.DonorID != donorID {
		return nil, errors.New("donation not found")
	}

	return donation, nil
}

// UpdateComment: the donor modifies their own comment.
// If the organizer had hidden it, it remains hidden.
func (s *DonationService) UpdateComment(ctx context.Context, donorID, donationID, comment string) error {
	if _, err := s.donationOwnedBy(ctx, donorID, donationID); err != nil {
		return err
	}

	cleaned, err := cleanComment(comment)
	if err != nil {
		return err
	}

	if err := s.donationRepo.UpdateComment(ctx, donationID, cleaned); err != nil {
		return fmt.Errorf("failed to update comment: %w", err)
	}
	return nil
}

// DeleteComment: the donor deletes their own comment.
func (s *DonationService) DeleteComment(ctx context.Context, donorID, donationID string) error {
	if _, err := s.donationOwnedBy(ctx, donorID, donationID); err != nil {
		return err
	}

	if err := s.donationRepo.UpdateComment(ctx, donationID, ""); err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}
	return nil
}

// SetCommentHidden: the organizer hides or unhides a comment on THEIR campaign.
// campaignID must have already been verified as belonging to the logged-in organizer.
func (s *DonationService) SetCommentHidden(ctx context.Context, campaignID, donationID string, hidden bool) error {
	if campaignID == "" || donationID == "" {
		return errors.New("campaign ID and donation ID are required")
	}

	donation, err := s.donationRepo.GetDonationByID(ctx, donationID)
	if err != nil {
		return fmt.Errorf("failed to fetch donation: %w", err)
	}
	if donation == nil || donation.CampaignID != campaignID {
		return errors.New("donation not found")
	}

	if err := s.donationRepo.SetCommentHidden(ctx, donationID, hidden); err != nil {
		return fmt.Errorf("failed to change comment visibility: %w", err)
	}
	return nil
}
