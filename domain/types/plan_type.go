package types

type PlanType string

const (
	Free PlanType = "free"
	Pro  PlanType = "pro"
)

const (
	FreeMonthlyLimit = 600  // 10 min in seconds
	ProMonthlyLimit  = 6000 // 100 min in seconds
)

func GetMonthlyLimit(plan PlanType) float64 {
	switch plan {
	case Pro:
		return ProMonthlyLimit
	default:
		return FreeMonthlyLimit
	}
}
