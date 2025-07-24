package config

type Env struct {
	AppEnv                 string `mapstructure:"APP_ENV"`
	ServerAddress          string `mapstructure:"SERVER_ADDRESS" validate:"required"`
	JWTSecret              string `mapstructure:"JWT_SECRET" validate:"required"`
	FrontEndURL            string `mapstructure:"FRONTEND_URL" validate:"required"`
	MongoURI               string `mapstructure:"MONGO_URI" validate:"required"`
	MongoDBName            string `mapstructure:"MONGO_DB_NAME" validate:"required"`
	GoogleRedirectURL      string `mapstructure:"GOOGLE_REDIRECT_URL" validate:"required"`
	GoogleClientID         string `mapstructure:"GOOGLE_CLIENT_ID" validate:"required"`
	GoogleClientSecret     string `mapstructure:"GOOGLE_CLIENT_SECRET" validate:"required"`
	GitHubRedirectURL      string `mapstructure:"GITHUB_REDIRECT_URL" validate:"required"`
	GitHubClientID         string `mapstructure:"GITHUB_CLIENT_ID" validate:"required"`
	GitHubClientSecret     string `mapstructure:"GITHUB_CLIENT_SECRET" validate:"required"`
	AWSRegion              string `mapstructure:"AWS_REGION" validate:"required"`
	AWSAccessKeyID         string `mapstructure:"AWS_ACCESS_KEY_ID" validate:"required"`
	AWSSecretAccessKey     string `mapstructure:"AWS_SECRET_ACCESS_KEY" validate:"required"`
	AWSS3BucketName        string `mapstructure:"AWS_S3_BUCKET_NAME" validate:"required"`
	AWSLambdaFuncName      string `mapstructure:"AWS_LAMBDA_FUNC_NAME" validate:"required"`
	SinchAppKey            string `mapstructure:"SINCH_APP_KEY" validate:"required"`
	SinchAppSecret         string `mapstructure:"SINCH_APP_SECRET" validate:"required"`
	ResendApiKey           string `mapstructure:"RESEND_API_KEY" validate:"required"`
	NotifyEmail            string `mapstructure:"NOTIFY_EMAIL" validate:"required"`
	PaddleAPIKey           string `mapstructure:"PADDLE_API_KEY" validate:"required"`
	PaddleWebhookSecretKey string `mapstructure:"PADDLE_WEBHOOK_SECRET_KEY" validate:"required"`
	SentryDSN              string `mapstructure:"SENTRY_DSN" validate:"required"`
}
