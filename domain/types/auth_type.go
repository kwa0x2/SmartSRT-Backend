package types

type AuthType string

const (
	Google      AuthType = "google"
	Github      AuthType = "github"
	Credentials AuthType = "credentials"
)
