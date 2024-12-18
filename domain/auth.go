package domain

type CredentialsSignUpBody struct {
	Name      string `json:"name"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	AvatarURL string `json:"avatar_url"`
}

type CredentialsSignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
