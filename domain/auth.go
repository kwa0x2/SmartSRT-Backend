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

type SinchSendOTPBody struct {
	PhoneNumber string `json:"phone_number"`
}

type SinchVerifyOTPBody struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}
