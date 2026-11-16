package models

import "time"

type DonationStatus string

const (
	DonationPending   DonationStatus = "pending"
	DonationCompleted DonationStatus = "completed"
	DonationFailed    DonationStatus = "failed"
	DonationRefunded  DonationStatus = "refunded"
)

type Donation struct {
	ID              string         `json:"id"`
	CampaignID      string         `json:"campaign_id"`
	PayoutAccountID string         `json:"payout_account_id"`
	DonorID         *string        `json:"donor_id,omitempty"`
	DonorName       string         `json:"donor_name,omitempty"`
	DonorPhone      string         `json:"donor_phone,omitempty"`
	DonorEmail      string         `json:"donor_email,omitempty"`
	Amount          int64          `json:"amount"`
	FeeAmount       int64          `json:"fee_amount"`    // platform fee, paid IN ADDITION by the donor
	TipAmount       int64          `json:"tip_amount"`    // optional tip
	TotalCharged    int64          `json:"total_charged"` // Amount + FeeAmount + TipAmount
	Currency        string         `json:"currency"`
	PaymentMethod   PaymentMethod  `json:"payment_method"`
	Status          DonationStatus `json:"status"`
	TransactionRef  string         `json:"transaction_ref,omitempty"`
	IsAnonymous     bool           `json:"is_anonymous"`
	Comment         string         `json:"comment,omitempty"`         // donor's comment (public once the donation is confirmed)
	CommentHidden   bool           `json:"comment_hidden,omitempty"`  // hidden by the campaign organizer
	DonorPhotoURL   string         `json:"donor_photo_url,omitempty"` // photo displayed with the donation (never if anonymous)
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	PayoutRef       string         `json:"payout_ref,omitempty"`
	Reply           string         `json:"reply,omitempty"` // organizer's reply to the comment
	ReplyAt         *time.Time     `json:"reply_at,omitempty"`
}
