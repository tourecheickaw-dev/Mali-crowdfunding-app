package models

import "time"

type PayoutType string

const (
	PayoutOrangeMoney PayoutType = "orange_money"
	PayoutWave        PayoutType = "wave"
	PayoutMoov        PayoutType = "moov_money"
	PayoutBank        PayoutType = "bank"
	PayoutMTN         PayoutType = "mtn_momo"
	PayoutFreeMoney   PayoutType = "free_money"
	PayoutWizall      PayoutType = "wizall"
)

type PayoutAccount struct {
	ID            string     `json:"id"`
	OrganizerID   string     `json:"organizer_id"`
	Type          PayoutType `json:"type"`
	AccountName   string     `json:"account_name"`        // Nom du titulaire tel qu'enregistré chez l'opérateur/banque
	AccountNumber string     `json:"account_number"`      // N° Téléphone (+223...) ou IBAN
	BankName      string     `json:"bank_name,omitempty"` // Requis si Type == "bank"
	BankCode      string     `json:"bank_code,omitempty"` // Code SWIFT/BIC ou RIB si banque
	Country       string     `json:"country"`
	IsActive      bool       `json:"is_active"`
	VerifiedAt    *time.Time `json:"verified_at,omitempty"` // Pointeur pour gérer le cas non encore vérifié (nil)
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}
