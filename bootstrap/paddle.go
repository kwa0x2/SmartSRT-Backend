package bootstrap

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"log"
)

func CreatePaddle(env *Env) *paddle.SDK {
	sdk, err := paddle.New(env.PADDLE_API_KEY, paddle.WithBaseURL(paddle.SandboxBaseURL))
	if err != nil {
		log.Fatal(err)
	}

	return sdk
}
