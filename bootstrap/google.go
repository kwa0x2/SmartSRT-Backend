package bootstrap

import (
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig(env *config.Env) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  env.GoogleRedirectURL,
		ClientID:     env.GoogleClientID,
		ClientSecret: env.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
