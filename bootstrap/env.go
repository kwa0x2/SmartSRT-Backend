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

	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logger.Warn("No .env file found, relying on environment variables",
			slog.String("file", ".env"),
			slog.String("error", err.Error()),
		)
	}

	viper.AutomaticEnv()

	viper.BindEnv("APP_ENV")
	viper.BindEnv("SERVER_ADDRESS")
	viper.BindEnv("JWT_SECRET")
	viper.BindEnv("FRONTEND_URL")
	viper.BindEnv("MONGO_URI")
	viper.BindEnv("MONGO_DB_NAME")
	viper.BindEnv("GOOGLE_REDIRECT_URL")
	viper.BindEnv("GOOGLE_CLIENT_ID")
	viper.BindEnv("GOOGLE_CLIENT_SECRET")
	viper.BindEnv("GITHUB_REDIRECT_URL")
	viper.BindEnv("GITHUB_CLIENT_ID")
	viper.BindEnv("GITHUB_CLIENT_SECRET")
	viper.BindEnv("AWS_REGION")
	viper.BindEnv("AWS_ACCESS_KEY_ID")
	viper.BindEnv("AWS_SECRET_ACCESS_KEY")
	viper.BindEnv("AWS_S3_BUCKET_NAME")
	viper.BindEnv("AWS_LAMBDA_FUNC_NAME")
	viper.BindEnv("SINCH_APP_KEY")
	viper.BindEnv("SINCH_APP_SECRET")
	viper.BindEnv("RESEND_API_KEY")
	viper.BindEnv("NOTIFY_EMAIL")
	viper.BindEnv("PADDLE_API_KEY")
	viper.BindEnv("PADDLE_WEBHOOK_SECRET_KEY")
	viper.BindEnv("SENTRY_DSN")
	viper.BindEnv("FREE_MONTHLY_LIMIT")
	viper.BindEnv("PRO_MONTHLY_LIMIT")

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
