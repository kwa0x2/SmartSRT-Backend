package bootstrap

import (
	"log/slog"
	"os"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func CreatePaddle(env *config.Env) *paddle.SDK {
	logger := slog.Default()

	sdk, err := paddle.New(env.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		logger.Error("Paddle SDK creation failed",
			slog.String("error", err.Error()),
			slog.String("base_url", paddle.SandboxBaseURL),
		)
		os.Exit(1)
	}

	logger.Info("Paddle SDK created successfully",
		slog.String("status", "initialized"),
	)

	return sdk
}
