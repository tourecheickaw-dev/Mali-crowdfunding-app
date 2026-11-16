package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

// donationInCampaign returns a donation only if it belongs to this campaign.
func (s *DonationService) donationInCampaign(ctx context.Context, campaignID, donationID string) (*models.Donation, error) {
	if campaignID == "" || donationID == "" {
		return nil, errors.New("campaign ID and donation ID are required")
	}

	donation, err := s.donationRepo.GetDonationByID(ctx, donationID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donation: %w", err)
	}
	if donation == nil || donation.CampaignID != campaignID {
		return nil, errors.New("donation not found")
	}

	return donation, nil
}

// ReplyToComment: the organizer replies to a donation comment on THEIR campaign.
// campaignID must already have been verified as belonging to the logged-in organizer.
func (s *DonationService) ReplyToComment(ctx context.Context, campaignID, donationID, reply string) error {
	donation, err := s.donationInCampaign(ctx, campaignID, donationID)
	if err != nil {
		return err
	}
	if donation.Status != models.DonationCompleted || donation.Comment == "" {
		return errors.New("there is no comment to reply to")
	}

	cleaned, err := cleanComment(reply)
	if err != nil {
		return err
	}
	if cleaned == "" {
		return errors.New("reply cannot be empty")
	}

	if err := s.donationRepo.SetReply(ctx, donationID, cleaned); err != nil {
		return fmt.Errorf("failed to save reply: %w", err)
	}
	return nil
}

// DeleteReply: the organizer removes their reply.
func (s *DonationService) DeleteReply(ctx context.Context, campaignID, donationID string) error {
	if _, err := s.donationInCampaign(ctx, campaignID, donationID); err != nil {
		return err
	}

	if err := s.donationRepo.SetReply(ctx, donationID, ""); err != nil {
		return fmt.Errorf("failed to delete reply: %w", err)
	}
	return nil
}
