package models

import "time"

type VerificationMethod string
type VerificationStatus string

const (
	DonorVerifOrangeOAuth VerificationMethod = "orange_money_oauth"
	DonorVerifWaveOAuth   VerificationMethod = "wave_oauth"
	DonorVerifSMS         VerificationMethod = "sms"
	DonorVerifEmail       VerificationMethod = "email"
	DonorVerifNone        VerificationMethod = "none"

	DonorUnverified VerificationStatus = "unverified"
	DonorVerified   VerificationStatus = "verified"
)

type Donor struct {
	ID                 string             `json:"id"`
	Phone              string             `json:"phone,omitempty"` // Ex: +22370000000
	Email              string             `json:"email,omitempty"`
	Name               string             `json:"name,omitempty"`
	ProfilePhotoURL    string             `json:"profile_photo_url,omitempty"`
	VerificationMethod VerificationMethod `json:"verification_method,omitempty"` // "sms", "email"
	VerificationStatus VerificationStatus `json:"verification_status"`           // "unverified", "verified"
	CreatedAt          time.Time          `json:"created_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
}
