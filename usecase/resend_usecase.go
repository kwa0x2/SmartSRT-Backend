package usecase

import (
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type resendUseCase struct {
	resendRepository domain.ResendRepository
}

func NewResendUseCase(resendRepository domain.ResendRepository) domain.ResendUseCase {
	return &resendUseCase{resendRepository: resendRepository}
}

func (ru *resendUseCase) SendSetupPasswordEmail(email, setupPassLink string) (string, error) {
	htmlContent, err := utils.LoadRecoveryEmailTemplate(setupPassLink)
	if err != nil {
		return "", err
	}

	sentID, err := ru.resendRepository.SendEmail(email, "set a new password", htmlContent)
	if err != nil {
		return "", err
	}
	return sentID, nil
}

func (ru *resendUseCase) SendContactNotifyMail(env *config.Env, contact *domain.Contact) (string, error) {
	htmlContent, err := utils.LoadContactNotifyTemplate(contact)
	if err != nil {
		return "", err
	}

	sentID, err := ru.resendRepository.SendEmail(env.NotifyEmail, "new contact form", htmlContent)
	if err != nil {
		return "", err
	}
	return sentID, nil
}

func (ru *resendUseCase) SendDeleteAccountEmail(email, deleteAccountLink string) (string, error) {
	htmlContent, err := utils.LoadDeleteAccountEmailTemplate(deleteAccountLink)
	if err != nil {
		return "", err
	}

	sentID, err := ru.resendRepository.SendEmail(email, "delete account", htmlContent)
	if err != nil {
		return "", err
	}
	return sentID, nil
}
