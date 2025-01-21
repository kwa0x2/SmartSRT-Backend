package usecase

import (
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type resendUseCase struct {
	resendRepository domain.ResendRepository
}

func NewResendUseCase(resendRepository domain.ResendRepository) domain.ResendUseCase {
	return &resendUseCase{resendRepository: resendRepository}
}

func (ru *resendUseCase) SendRecoveryEmail(email, recoveryLink string) (string, error) {
	htmlContent, err := utils.LoadRecoveryEmailTemplate(email, recoveryLink)
	if err != nil {
		return "", err
	}

	sentID, err := ru.resendRepository.SendEmail(email, "reset password", htmlContent)
	if err != nil {
		return "", err
	}
	return sentID, nil
}
