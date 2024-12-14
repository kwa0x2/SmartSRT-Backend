package bootstrap

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

func GitHubConfig(env *Env) oauth2.Config {
	return oauth2.Config{
		ClientID:     env.GitHubClientID,
		ClientSecret: env.GitHubClientSecret,
		RedirectURL:  env.GitHubRedirectURL,
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
}
