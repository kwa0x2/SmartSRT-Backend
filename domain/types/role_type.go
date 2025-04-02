package types

type RoleType string

const (
	Free RoleType = "free"
	Pro  RoleType = "pro"
)

const (
	FreeMonthlyLimit = 600  // 10 min in seconds
	ProMonthlyLimit  = 6000 // 100 min in seconds
)

func GetMonthlyLimit(role RoleType) float64 {
	switch role {
	case Pro:
		return ProMonthlyLimit
	default:
		return FreeMonthlyLimit
	}
}
