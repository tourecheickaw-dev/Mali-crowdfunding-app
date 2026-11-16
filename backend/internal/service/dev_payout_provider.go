package service

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
)

// DevPayoutProvider simulates an aggregator: sends NO money. Local testing only.
type DevPayoutProvider struct{}

func (DevPayoutProvider) SendPayout(_ context.Context, req PayoutRequest) (PayoutResult, error) {
	log.Printf("[SIMULATION] transfer of %d %s via %s, reference %s (no money sent)",
		req.Amount, req.Currency, req.Account.Type, req.Reference)
	return PayoutResult{TransferID: "dev-" + uuid.New().String()}, nil
}

// NoPayoutProvider is used when no aggregator is configured: it rejects any transfer
// instead of pretending.
type NoPayoutProvider struct{}

func (NoPayoutProvider) SendPayout(_ context.Context, _ PayoutRequest) (PayoutResult, error) {
	return PayoutResult{}, errors.New("no payment provider configured")
}
