package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/repository"
)

type DonationService struct {
	donationRepo  *repository.DonationRepository
	campaignRepo  *repository.CampaignRepository
	payoutRepo    *repository.PayoutAccountRepository
	organiserRepo *repository.OrganizerRepository
	provider      PayoutProvider
}

func NewDonationService(
	donationRepo *repository.DonationRepository,
	campaignRepo *repository.CampaignRepository,
	payoutRepo *repository.PayoutAccountRepository,
	organiserRepo *repository.OrganizerRepository,
	provider PayoutProvider,
) *DonationService {
	return &DonationService{
		donationRepo:  donationRepo,
		campaignRepo:  campaignRepo,
		payoutRepo:    payoutRepo,
		organiserRepo: organiserRepo,
		provider:      provider,
	}
}

type CreateDonationInput struct {
	CampaignID    string
	DonorID       *string
	DonorName     string
	DonorPhone    string
	DonorEmail    string
	Amount        int64  // amount received by the organizer (excluding fee and tip)
	TipAmount     int64  // optional tip for the platform
	Comment       string // optional donor comment
	DonorPhotoURL string // photo displayed with the donation (never for an anonymous donation)
	PaymentMethod models.PaymentMethod
	IsAnonymous   bool
}

// CreateDonation records a pending donation with the organizer's active account.
func (s *DonationService) CreateDonation(ctx context.Context, input CreateDonationInput) (*models.Donation, error) {
	if input.CampaignID == "" {
		return nil, errors.New("campaign ID is required")
	}
	if input.Amount <= 0 {
		return nil, errors.New("amount must be greater than zero")
	}
	if input.TipAmount < 0 {
		return nil, errors.New("tip amount cannot be negative")
	}
	if input.TipAmount > input.Amount {
		return nil, errors.New("tip amount cannot be greater than the donation")
	}

	comment, err := cleanComment(input.Comment)
	if err != nil {
		return nil, err
	}

	campaign, err := s.campaignRepo.GetCampaignByID(ctx, input.CampaignID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch campaign: %w", err)
	}
	if campaign == nil {
		return nil, errors.New("campaign not found")
	}
	if campaign.Status != models.CampaignActive {
		return nil, errors.New("campaign is not accepting donations")
	}

	account, err := s.activePayoutAccount(ctx, campaign.OrganizerID)
	if err != nil {
		return nil, err
	}

	// Never show photo on an anonymous donation.
	donorPhoto := input.DonorPhotoURL
	if input.IsAnonymous {
		donorPhoto = ""
	}

	fee := CalculateFee(input.Amount)
	now := time.Now()
	donation := models.Donation{
		ID:              uuid.New().String(),
		CampaignID:      campaign.ID,
		PayoutAccountID: account.ID,
		DonorID:         input.DonorID,
		DonorName:       input.DonorName,
		DonorPhone:      input.DonorPhone,
		DonorEmail:      input.DonorEmail,
		Amount:          input.Amount,
		FeeAmount:       fee,
		TipAmount:       input.TipAmount,
		TotalCharged:    input.Amount + fee + input.TipAmount,
		Comment:         comment,
		DonorPhotoURL:   donorPhoto,
		Currency:        campaign.Currency,
		PaymentMethod:   input.PaymentMethod,
		Status:          models.DonationPending,
		IsAnonymous:     input.IsAnonymous,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.donationRepo.CreateDonation(ctx, &donation); err != nil {
		return nil, fmt.Errorf("failed to save donation: %w", err)
	}

	return &donation, nil
}

// CompleteDonation is called when the aggregator confirms the donor's payment.
// It can be called multiple times for the same donation (replayed webhook) without duplicate effects.
func (s *DonationService) CompleteDonation(ctx context.Context, donationID string) error {
	donation, err := s.donationRepo.GetDonationByID(ctx, donationID)
	if err != nil {
		return fmt.Errorf("failed to fetch donation: %w", err)
	}
	if donation == nil {
		return errors.New("donation not found")
	}
	if donation.Status == models.DonationFailed || donation.Status == models.DonationRefunded {
		return fmt.Errorf("donation %s is %s, cannot complete it", donation.ID, donation.Status)
	}

	// 1. Count the donation (only once)
	if donation.Status == models.DonationPending {
		donation.Status = models.DonationCompleted
		donation.UpdatedAt = time.Now()
		if err := s.donationRepo.CreateDonation(ctx, donation); err != nil {
			return fmt.Errorf("failed to mark donation as completed: %w", err)
		}
		if err := s.campaignRepo.UpdateCampaignAmount(ctx, donation.CampaignID, donation.Amount); err != nil {
			return fmt.Errorf("failed to update campaign amount: %w", err)
		}
		s.addToOrganizerTotal(ctx, donation)
	}

	// 2. Send the money to the organizer, only if it hasn't been sent already
	if donation.PayoutRef == "" {
		return s.sendPayout(ctx, donation)
	}
	return nil
}

// addToOrganizerTotal updates the organizer's lifetime total raised.
// Best effort: an error is written to the logs but never blocks sending the money.
func (s *DonationService) addToOrganizerTotal(ctx context.Context, donation *models.Donation) {
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, donation.CampaignID)
	if err != nil || campaign == nil {
		log.Printf("organizer total not updated (donation %s): campaign not found (%v)", donation.ID, err)
		return
	}
	if err := s.organiserRepo.AddToTotalRaised(ctx, campaign.OrganizerID, donation.Amount); err != nil {
		log.Printf("organizer total not updated (donation %s): %v", donation.ID, err)
	}
}

// sendPayout sends the organizer's share to their payout account.
// If this account refuses, switch once to a backup account.
func (s *DonationService) sendPayout(ctx context.Context, donation *models.Donation) error {
	campaign, err := s.campaignRepo.GetCampaignByID(ctx, donation.CampaignID)
	if err != nil {
		return fmt.Errorf("failed to fetch campaign: %w", err)
	}
	if campaign == nil {
		return errors.New("campaign not found")
	}

	accounts, err := s.payoutRepo.GetPayoutAccountsByOrganizerID(ctx, campaign.OrganizerID)
	if err != nil {
		return fmt.Errorf("failed to fetch payout accounts: %w", err)
	}

	current := findAccount(accounts, donation.PayoutAccountID)
	if current == nil {
		return errors.New("payout account of this donation not found")
	}

	result, err := s.provider.SendPayout(ctx, PayoutRequest{
		Reference: donation.ID,
		Amount:    donation.Amount,
		Currency:  donation.Currency,
		Account:   *current,
	})
	if err != nil {
		// Network error or timeout: the transfer might have gone through. Do not resend without checking its status.
		return fmt.Errorf("payout request failed, check its status before retrying: %w", err)
	}

	switch result.Failure {
	case "":
		return s.recordPayout(ctx, donation, result.TransferID)
	case FailureDestination:
		return s.payoutToFallback(ctx, donation, campaign.OrganizerID, accounts)
	case FailureTemporary:
		return fmt.Errorf("temporary payout failure for donation %s, retry later", donation.ID)
	default:
		return fmt.Errorf("payout for donation %s needs manual review", donation.ID)
	}
}

// payoutToFallback attempts the transfer once on another account of the organizer.
func (s *DonationService) payoutToFallback(ctx context.Context, donation *models.Donation, organizerID string, accounts []models.PayoutAccount) error {
	var fallback *models.PayoutAccount
	for i := range accounts {
		if accounts[i].ID != donation.PayoutAccountID {
			fallback = &accounts[i]
			break
		}
	}
	if fallback == nil {
		return fmt.Errorf("payout for donation %s refused and no fallback account available", donation.ID)
	}

	result, err := s.provider.SendPayout(ctx, PayoutRequest{
		Reference: donation.ID + "-fb", // different reference from the 1st attempt
		Amount:    donation.Amount,
		Currency:  donation.Currency,
		Account:   *fallback,
	})
	if err != nil {
		return fmt.Errorf("fallback payout request failed, check its status before retrying: %w", err)
	}
	if result.Failure != "" {
		return fmt.Errorf("fallback payout for donation %s failed (%s), manual review needed", donation.ID, result.Failure)
	}

	// The money has been sent: record it first, then switch the active account.
	donation.PayoutAccountID = fallback.ID
	if err := s.recordPayout(ctx, donation, result.TransferID); err != nil {
		return err
	}
	if err := s.payoutRepo.SetActivePayoutAccount(ctx, organizerID, fallback.ID); err != nil {
		return fmt.Errorf("payout sent but failed to switch active account: %w", err)
	}
	return nil
}

// recordPayout records the transfer identifier on the donation.
func (s *DonationService) recordPayout(ctx context.Context, donation *models.Donation, transferID string) error {
	donation.PayoutRef = transferID
	donation.UpdatedAt = time.Now()
	if err := s.donationRepo.CreateDonation(ctx, donation); err != nil {
		return fmt.Errorf("payout sent (%s) but failed to record it: %w", transferID, err)
	}
	return nil
}

func (s *DonationService) activePayoutAccount(ctx context.Context, organizerID string) (*models.PayoutAccount, error) {
	accounts, err := s.payoutRepo.GetPayoutAccountsByOrganizerID(ctx, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payout accounts: %w", err)
	}
	for i := range accounts {
		if accounts[i].IsActive {
			return &accounts[i], nil
		}
	}
	return nil, errors.New("organizer has no active payout account")
}

func findAccount(accounts []models.PayoutAccount, id string) *models.PayoutAccount {
	for i := range accounts {
		if accounts[i].ID == id {
			return &accounts[i]
		}
	}
	return nil
}
