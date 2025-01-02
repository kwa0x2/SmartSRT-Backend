package domain

type CredentialsSignUpBody struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	AvatarURL   string `json:"avatar_url"`
	OTP         string `json:"otp"`
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
