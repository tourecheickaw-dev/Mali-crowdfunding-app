package service

const FeePercent = 6

// ComputeFee returns the fee in FCFA, rounded to the nearest franc.
func ComputeFee(amount int64) int64 {
	if amount <= 0 {
		return 0
	}
	return (amount*FeePercent + 50) / 100
}
