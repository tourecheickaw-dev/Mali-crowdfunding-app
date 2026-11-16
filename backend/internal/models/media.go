package models

import "time"

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
)

type CampaignMedia struct {
	ID           string    `json:"id"`
	CampaignID   string    `json:"campaign_id"`
	URL          string    `json:"url"`
	MediaType    MediaType `json:"media_type"`              // "image" ou "video"
	Caption      string    `json:"caption,omitempty"`       // Légende/Description de la photo
	IsCover      bool      `json:"is_cover"`                // Image principale
	IsSensitive  bool      `json:"is_sensitive"`            // Pour avertissement contenu sensible
	DisplayOrder int       `json:"display_order,omitempty"` // Pour ordonner la galerie
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
