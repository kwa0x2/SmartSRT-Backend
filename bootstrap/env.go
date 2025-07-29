package bootstrap

import (
	"log/slog"
	"os"

	"github.com/kwa0x2/AutoSRT-Backend/config"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewEnv() *config.Env {
	logger := slog.Default()
	env := config.Env{}

	viper.SetConfigFile(".env") // pls run main.go in cmd dir

	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("Environment file could not be read",
			slog.String("file", ".env"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		logger.Error("Environment file could not be parsed",
			slog.String("file", ".env"),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	validate := validator.New()
	if err = validate.Struct(env); err != nil {
		logger.Error("Environment validation failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("Environment loaded successfully",
		slog.String("app_env", env.AppEnv),
		slog.String("server_address", env.ServerAddress),
		slog.String("mongo_db_name", env.MongoDBName),
	)

	if env.AppEnv == "development" {
		logger.Debug("Application running in development mode",
			slog.String("mode", "development"),
		)
	}

	return &env
}
