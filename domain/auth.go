package domain

import "github.com/kwa0x2/AutoSRT-Backend/domain/types"

type CredentialsLoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SinchVerifyOTPBody struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type PasswordBody struct {
	Password string `json:"password"`
}

type VerifyOTPAndCreateBody struct {
	Name        string         `json:"name"`
	Email       string         `json:"email"`
	PhoneNumber string         `json:"phone_number"`
	Password    string         `json:"password"`
	AvatarURL   string         `json:"avatar_url"`
	OTP         string         `json:"otp"`
	AuthType    types.AuthType `json:"auth_type"`
}

type EmailBody struct {
	Email string `json:"email"`
}

type PhoneNumberBody struct {
	PhoneNumber string `json:"phone_number"`
}
