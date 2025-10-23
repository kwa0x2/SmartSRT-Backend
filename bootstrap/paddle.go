package bootstrap

import (
	"log/slog"
	"os"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/SmartSRT-Backend/config"
)

func CreatePaddle(env *config.Env) *paddle.SDK {
	logger := slog.Default()

	baseURL := paddle.SandboxBaseURL
	if env.AppEnv == "production" {
		baseURL = paddle.ProductionBaseURL
	}

	sdk, err := paddle.New(env.PaddleAPIKey, paddle.WithBaseURL(baseURL))
	if err != nil {
		logger.Error("Paddle SDK creation failed",
			slog.String("error", err.Error()),
			slog.String("base_url", baseURL),
		)
		os.Exit(1)
	}

	logger.Info("Paddle SDK created successfully",
		slog.String("status", "initialized"),
		slog.String("base_url", baseURL),
	)

	return sdk
}
