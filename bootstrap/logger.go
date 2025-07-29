package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func SetupLogger(env *config.Env) *slog.Logger {
	var level slog.Level

	switch env.AppEnv {
	case "production":
		level = slog.LevelInfo
	case "development":
		level = slog.LevelDebug
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)

	return logger
}

func LogWithFields(logger *slog.Logger, level slog.Level, msg string, fields map[string]interface{}) {
	args := make([]interface{}, 0, len(fields)*2)
	for key, value := range fields {
		args = append(args, key, value)
	}

	logger.Log(context.TODO(), level, msg, args...)
}
