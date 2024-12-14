package bootstrap

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func GoogleConfig(env *Env) oauth2.Config {
	return oauth2.Config{
		RedirectURL:  env.GoogleRedirectURL,
		ClientID:     env.GoogleClientID,
		ClientSecret: env.GoogleClientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}
