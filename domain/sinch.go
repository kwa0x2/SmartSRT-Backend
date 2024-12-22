package domain

type SinchRepository interface {
	SendOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, code string) (bool, error)
}

type SinchUseCase interface {
	SendOTP(phoneNumber string) error
	VerifyOTP(phoneNumber, code string) (bool, error)
}
