package usecase

import "github.com/kwa0x2/AutoSRT-Backend/domain"

type sinchUseCase struct {
	sinchRepository domain.SinchRepository
}

func NewSinchUseCase(sinchRepository domain.SinchRepository) domain.SinchUseCase {
	return &sinchUseCase{
		sinchRepository: sinchRepository,
	}
}

func (su *sinchUseCase) SendOTP(phoneNumber string) error {
	return su.sinchRepository.SendOTP(phoneNumber)
}

func (su *sinchUseCase) VerifyOTP(phoneNumber, code string) (bool, error) {
	return su.sinchRepository.VerifyOTP(phoneNumber, code)
}
