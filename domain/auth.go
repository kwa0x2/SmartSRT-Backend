package domain

import "github.com/kwa0x2/AutoSRT-Backend/domain/types"

type VerifyOTPAndCreateBody struct {
	Name        string             `json:"name"`
	Email       string             `json:"email"`
	PhoneNumber string             `json:"phone_number"`
	Password    string             `json:"password"`
	AvatarURL   string             `json:"avatar_url"`
	OTP         string             `json:"otp"`
	AuthWith    types.AutoWithType `json:"auth_with"`
}

type CredentialsSignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SinchVerifyOTPBody struct {
	PhoneNumber string `json:"phone_number"`
	Code        string `json:"code"`
}

type IsEmailExistsBody struct {
	Email string `json:"email"`
}

type PhoneNumberBody struct {
	PhoneNumber string `json:"phone_number"`
}
