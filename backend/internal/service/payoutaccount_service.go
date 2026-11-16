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

type PayoutAccountService struct {
	payoutRepo *repository.PayoutAccountRepository
}

func NewPayoutAccountService(payoutRepo *repository.PayoutAccountRepository) *PayoutAccountService {
	return &PayoutAccountService{
		payoutRepo: payoutRepo,
	}
}

type AddPayoutAccountInput struct {
	OrganizerID   string
	Type          models.PayoutType
	AccountName   string // Account holder name as registered with the operator/bank
	AccountNumber string // Phone number (+country code) for mobile money, IBAN/RIB for a bank
	BankName      string // Required for a bank account
	BankCode      string // Optional
	Country       string // Required for a bank account (ISO code, e.g. "ML"); derived from phone for mobile money
}

// AddPayoutAccount adds a payout account for an organizer.
// The first added account automatically becomes the active account.
func (s *PayoutAccountService) AddPayoutAccount(ctx context.Context, input AddPayoutAccountInput) (*models.PayoutAccount, error) {
	if input.OrganizerID == "" {
		return nil, errors.New("organizer ID is required")
	}
	if input.AccountName == "" {
		return nil, errors.New("account holder name is required")
	}
	if input.AccountNumber == "" {
		return nil, errors.New("account number is required")
	}

	accountNumber := input.AccountNumber
	country := input.Country

	switch input.Type {
	case models.PayoutOrangeMoney, models.PayoutWave, models.PayoutMoov,
		models.PayoutMTN, models.PayoutFreeMoney, models.PayoutWizall:
		phone, err := normalizePhone(input.AccountNumber)
		if err != nil {
			return nil, fmt.Errorf("invalid mobile money number: %w", err)
		}
		if !isUEMOAPhone(phone) {
			return nil, errors.New("mobile money number must belong to the UEMOA zone")
		}
		accountNumber = phone
		country = countryFromPhone(phone)
	case models.PayoutBank:
		if input.BankName == "" {
			return nil, errors.New("bank name is required for a bank account")
		}
		if country == "" {
			return nil, errors.New("country is required for a bank account")
		}
	default:
		return nil, errors.New("payout type must be orange_money, wave, moov_money or bank")
	}

	existing, err := s.payoutRepo.GetPayoutAccountsByOrganizerID(ctx, input.OrganizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing payout accounts: %w", err)
	}
	for _, acc := range existing {
		if acc.Type == input.Type && acc.AccountNumber == accountNumber {
			return nil, errors.New("this payout account is already registered")
		}
	}

	now := time.Now()
	account := models.PayoutAccount{
		ID:            uuid.New().String(),
		OrganizerID:   input.OrganizerID,
		Type:          input.Type,
		AccountName:   input.AccountName,
		AccountNumber: accountNumber,
		BankName:      input.BankName,
		BankCode:      input.BankCode,
		Country:       country,
		IsActive:      len(existing) == 0, // First account is active by default
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.payoutRepo.SavePayoutAccount(ctx, &account); err != nil {
		return nil, fmt.Errorf("failed to save payout account: %w", err)
	}

	return &account, nil
}

// GetOrganizerPayoutAccounts retrieves all payout accounts associated with an organizer.
func (s *PayoutAccountService) GetOrganizerPayoutAccounts(ctx context.Context, organizerID string) ([]models.PayoutAccount, error) {
	if organizerID == "" {
		return nil, errors.New("organizer ID is required")
	}

	accounts, err := s.payoutRepo.GetPayoutAccountsByOrganizerID(ctx, organizerID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch payout accounts for organizer %s: %w", organizerID, err)
	}

	return accounts, nil
}

// SetActiveAccount changes the active account of an organizer (manual choice).
func (s *PayoutAccountService) SetActiveAccount(ctx context.Context, organizerID string, accountID string) error {
	if organizerID == "" {
		return errors.New("organizer ID is required")
	}
	if accountID == "" {
		return errors.New("account ID is required")
	}

	accounts, err := s.payoutRepo.GetPayoutAccountsByOrganizerID(ctx, organizerID)
	if err != nil {
		return fmt.Errorf("failed to fetch organizer accounts: %w", err)
	}

	found := false
	for _, acc := range accounts {
		if acc.ID == accountID {
			found = true
			break
		}
	}
	if !found {
		return errors.New("specified account does not belong to this organizer or does not exist")
	}

	if err := s.payoutRepo.SetActivePayoutAccount(ctx, organizerID, accountID); err != nil {
		return fmt.Errorf("failed to set active payout account: %w", err)
	}

	return nil
}
