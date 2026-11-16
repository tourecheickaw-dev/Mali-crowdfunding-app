package models

import "time"

type OrganizerVerificationStatus string

const (
	OrganizerUnverified OrganizerVerificationStatus = "unverified"
	OrganizerPending    OrganizerVerificationStatus = "pending"
	OrganizerVerified   OrganizerVerificationStatus = "verified"
	OrganizerRejected   OrganizerVerificationStatus = "rejected"
)

type Organizer struct {
	ID                  string                      `json:"id"`
	DisplayName         string                      `json:"display_name"`
	Phone               string                      `json:"phone"`
	Email               string                      `json:"email,omitempty"`
	Bio                 string                      `json:"bio,omitempty"`
	ProfilePhotoURL     string                      `json:"profile_photo_url,omitempty"`
	Country             string                      `json:"country"`
	VerificationStatus  OrganizerVerificationStatus `json:"verification_status"`
	TotalRaisedLifetime int64                       `json:"total_raised_lifetime"`
	CreatedAt           time.Time                   `json:"created_at"`
	UpdatedAt           time.Time                   `json:"updated_at"`
}
