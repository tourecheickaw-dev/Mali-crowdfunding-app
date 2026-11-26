package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

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

// accentReplacer removes French accents so "Touré" and "TOURE" are compared as equal.
var accentReplacer = strings.NewReplacer(
	"à", "a", "â", "a", "ä", "a", "á", "a", "ã", "a",
	"é", "e", "è", "e", "ê", "e", "ë", "e",
	"î", "i", "ï", "i", "í", "i",
	"ô", "o", "ö", "o", "ó", "o", "õ", "o",
	"ù", "u", "û", "u", "ü", "u", "ú", "u",
	"ç", "c", "ñ", "n", "œ", "oe", "æ", "ae",
)

// nameTokens splits a name into lowercase words without accents or punctuation.
func nameTokens(name string) []string {
	s := accentReplacer.Replace(strings.ToLower(name))
	return strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) })
}

// namesMatch compares two names ignoring case, accents, punctuation and word order.
// Every word of the shorter name must appear in the longer one
// (so a missing middle name is fine). A one-word name is only accepted
// if both names are exactly that word.
func namesMatch(a, b string) bool {
	ta, tb := nameTokens(a), nameTokens(b)
	if len(ta) == 0 || len(tb) == 0 {
		return false
	}
	if len(ta) > len(tb) {
		ta, tb = tb, ta // ta is now the shorter one
	}
	set := make(map[string]bool, len(tb))
	for _, t := range tb {
		set[t] = true
	}
	for _, t := range ta {
		if !set[t] {
			return false
		}
	}
	if len(ta) < 2 && len(tb) != len(ta) {
		return false
	}
	return true
}

func isMobileMoney(t models.PayoutType) bool {
	switch t {
	case models.PayoutOrangeMoney, models.PayoutWave, models.PayoutMoov,
		models.PayoutMTN, models.PayoutFreeMoney, models.PayoutWizall:
		return true
	}
	return false
}

// checkBankHolder accepts a bank account only if the organizer already has a verified
// mobile money account whose holder has the same name.
// For local tests only, ALLOW_UNVERIFIED_MOBILE_MONEY=true accepts unverified mobile money accounts.
func checkBankHolder(existing []models.PayoutAccount, bankHolder string) error {
	allowUnverified := os.Getenv("ALLOW_UNVERIFIED_MOBILE_MONEY") == "true"
	hasMobile := false
	for _, acc := range existing {
		if !isMobileMoney(acc.Type) {
			continue
		}
		if acc.VerifiedAt == nil && !allowUnverified {
			continue
		}
		hasMobile = true
		if namesMatch(acc.AccountName, bankHolder) {
			return nil
		}
	}
	if !hasMobile {
		return errors.New("add a verified mobile money account first")
	}
	return errors.New("the bank account holder name must match the name of your verified mobile money account")
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
	var verifiedAt *time.Time

	// A bank account is accepted only if it matches a verified mobile money account.
	if input.Type == models.PayoutBank {
		if err := checkBankHolder(existing, input.AccountName); err != nil {
			return nil, err
		}
		verifiedAt = &now
	}

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
		VerifiedAt:    verifiedAt,
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
