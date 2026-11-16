package service

// PlatformFeePercent: platform fee, added to the donation amount
// (the donor pays the donation + the fee, the organizer receives the full amount).
const PlatformFeePercent = 6

// CalculateFee returns the fee in whole FCFA, rounded to the nearest.
func CalculateFee(amount int64) int64 {
	return (amount*PlatformFeePercent + 50) / 100
}
