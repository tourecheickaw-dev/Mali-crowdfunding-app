package models

import "time"

type BeneficiaryType string
type CampaignStatus string

const (
	BeneficiarySelf            BeneficiaryType = "self"
	BeneficiaryOtherPerson     BeneficiaryType = "other_person"
	BeneficiaryCauseCollective BeneficiaryType = "cause_collective"

	CampaignDraft     CampaignStatus = "draft"
	CampaignActive    CampaignStatus = "active"
	CampaignCompleted CampaignStatus = "completed"
)

type Campaign struct {
	ID              string          `json:"id"`
	OrganizerID     string          `json:"organizer_id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	GoalAmount      *int64          `json:"goal_amount,omitempty"`
	CurrentAmount   int64           `json:"current_amount"`
	Currency        string          `json:"currency"`
	BeneficiaryType BeneficiaryType `json:"beneficiary_type"`
	BeneficiaryName string          `json:"beneficiary_name,omitempty"`
	Country         string          `json:"country"`
	URLSlug         string          `json:"url_slug"`
	Status          CampaignStatus  `json:"status"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}
