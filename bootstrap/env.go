package bootstrap

import (
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewEnv() *config.Env {
	env := config.Env{}
	viper.SetConfigFile(".env") // pls run main.go in cmd dir

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
