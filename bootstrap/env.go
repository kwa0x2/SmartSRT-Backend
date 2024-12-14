package bootstrap

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

type Env struct {
	AppEnv             string `mapstructure:"APP_ENV"`
	ServerAddress      string `mapstructure:"SERVER_ADDRESS" validate:"required"`
	MongoURI           string `mapstructure:"MONGO_URI" validate:"required"`
	MongoDBName        string `mapstructure:"MONGO_DB_NAME" validate:"required"`
	GoogleRedirectURL  string `mapstructure:"GOOGLE_REDIRECT_URL" validate:"required"`
	GoogleClientID     string `mapstructure:"GOOGLE_CLIENT_ID" validate:"required"`
	GoogleClientSecret string `mapstructure:"GOOGLE_CLIENT_SECRET" validate:"required"`
	GitHubRedirectURL  string `mapstructure:"GITHUB_REDIRECT_URL" validate:"required"`
	GitHubClientID     string `mapstructure:"GITHUB_CLIENT_ID" validate:"required"`
	GitHubClientSecret string `mapstructure:"GITHUB_CLIENT_SECRET" validate:"required"`
}

func NewEnv() *Env {
	env := Env{}
	viper.SetConfigFile("../.env") // pls run main.go in cmd dir

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Can't read .env file")
	}

	err = viper.Unmarshal(&env)
	if err != nil {
		log.Fatal("Can't unmarshal .env file")
	}

	validate := validator.New()
	if err = validate.Struct(env); err != nil {
		log.Fatal(err)
	}

	if env.AppEnv == "development" {
		log.Println("The App is running in development env")
	}

	return &env
}
