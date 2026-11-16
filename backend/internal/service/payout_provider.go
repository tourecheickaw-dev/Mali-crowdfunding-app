package service

import (
	"context"

	"github.com/tourecheickaw-dev/Mali-crowdfunding-app/internal/models"
)

// PayoutFailureKind classifies a transfer failure, regardless of the aggregator used.
type PayoutFailureKind string

const (
	// FailureDestination: recipient account is refused (-> fallback account).
	FailureDestination PayoutFailureKind = "destination"
	// FailureTemporary: temporary issue (-> retry later, same account).
	FailureTemporary PayoutFailureKind = "temporary"
	// FailureManual: check manually (insufficient funds, fraud, unknown reason...).
	FailureManual PayoutFailureKind = "manual"
)

// PayoutRequest describes a transfer to the organizer.
type PayoutRequest struct {
	Reference string
	Amount    int64
	Currency  string
	Account   models.PayoutAccount
}

// PayoutResult is the aggregator's response to sending a transfer.
type PayoutResult struct {
	TransferID string
	Failure    PayoutFailureKind // empty = transfer accepted
}

// PayoutProvider is what DonationService expects from an aggregator.
type PayoutProvider interface {
	SendPayout(ctx context.Context, req PayoutRequest) (PayoutResult, error)
}
