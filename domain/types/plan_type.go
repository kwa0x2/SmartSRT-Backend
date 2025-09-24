package types

import "github.com/kwa0x2/SmartSRT-Backend/config"

type PlanType string

const (
	Free PlanType = "free"
	Pro  PlanType = "pro"
)

func GetMonthlyLimit(plan PlanType, env *config.Env) float64 {
	switch plan {
	case Pro:
		return env.ProMonthlyLimit
	default:
		return env.FreeMonthlyLimit
	}
}
