package models

type PaymentMethod string

const (
	MethodOrangeMoney PaymentMethod = "orange_money"
	MethodWave        PaymentMethod = "wave"
	MethodMoov        PaymentMethod = "moov"
	MethodMTN         PaymentMethod = "mtn_momo"
	MethodFreeMoney   PaymentMethod = "free_money"
	MethodWizall      PaymentMethod = "wizall"
	MethodCard        PaymentMethod = "card"
)
