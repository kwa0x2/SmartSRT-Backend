package bootstrap

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"log"
)

func CreatePaddle(env *Env) *paddle.SDK {
	sdk, err := paddle.New(env.PaddleAPIKey, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		log.Fatal(err)
	}

	return sdk
}
