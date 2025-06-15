package bootstrap

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"log"
)

func CreatePaddle(env *config.Env) *paddle.SDK {
	sdk, err := paddle.New(env.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		log.Fatal(err)
	}

	return sdk
}
