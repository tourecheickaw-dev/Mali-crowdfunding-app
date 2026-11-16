package service

import (
	"errors"
	"strings"
)

// UEMOA zone country dialing codes.
var uemoaDialCodes = []string{
	"+229", // Benin
	"+226", // Burkina Faso
	"+225", // Cote d'Ivoire
	"+223", // Mali
	"+227", // Niger
	"+221", // Senegal
	"+228", // Togo
}

// normalizePhone cleans a phone number and enforces the international format (+dialing_code...).
// Accepts "00223..." and converts it to "+223...".
func normalizePhone(phone string) (string, error) {
	cleaned := strings.NewReplacer(" ", "", "-", "", ".", "", "(", "", ")", "").Replace(strings.TrimSpace(phone))

	if strings.HasPrefix(cleaned, "00") {
		cleaned = "+" + cleaned[2:]
	}
	if !strings.HasPrefix(cleaned, "+") {
		return "", errors.New("phone number must include the country code, e.g. +223...")
	}

	digits := cleaned[1:]
	if len(digits) < 8 || len(digits) > 15 {
		return "", errors.New("invalid phone number length")
	}
	for _, r := range digits {
		if r < '0' || r > '9' {
			return "", errors.New("phone number must contain only digits after the country code")
		}
	}

	return cleaned, nil
}

// isUEMOAPhone indicates whether an (already normalized) phone number belongs to the UEMOA zone.
// To be used for organizers and payout accounts, not for donors.
func isUEMOAPhone(phone string) bool {
	for _, code := range uemoaDialCodes {
		if strings.HasPrefix(phone, code) {
			return true
		}
	}
	return false
}

// Dialing code to ISO country code mapping for the UEMOA zone.
var uemoaCountryByDialCode = map[string]string{
	"+229": "BJ",
	"+226": "BF",
	"+225": "CI",
	"+223": "ML",
	"+227": "NE",
	"+221": "SN",
	"+228": "TG",
}

// countryFromPhone returns the ISO country code of a normalized UEMOA phone number ("" if unknown).
func countryFromPhone(phone string) string {
	for code, country := range uemoaCountryByDialCode {
		if strings.HasPrefix(phone, code) {
			return country
		}
	}
	return ""
}
