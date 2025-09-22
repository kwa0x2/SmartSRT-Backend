package bootstrap

import (
	"log/slog"
	"os"
	"strings"

	"github.com/kwa0x2/AutoSRT-Backend/config"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewEnv() *config.Env {
	logger := slog.Default()
	env := config.Env{}
	
	
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("No .env file found, relying on environment variables",
			slog.String("file", ".env"),
			slog.String("error", err.Error()),
		)
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&env); err != nil {
		logger.Error("Environment could not be parseds",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	validate := validator.New()
	if err := validate.Struct(env); err != nil {
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
