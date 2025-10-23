package usecase

import (
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
)

type resendUseCase struct {
	resendRepository domain.ResendRepository
}

func NewResendUseCase(resendRepository domain.ResendRepository) domain.ResendUseCase {
	return &resendUseCase{resendRepository: resendRepository}
}

func (ru *resendUseCase) sendEmail(to, subject string, templateLoader func() (string, error)) (string, error) {
	htmlContent, err := templateLoader()
	if err != nil {
		return "", err
	}

	return ru.resendRepository.SendEmail(to, subject, htmlContent)
}

func (ru *resendUseCase) SendSetupPasswordEmail(email, setupPassLink string) (string, error) {
	return ru.sendEmail(email, "🔐 SmartSRT - Set Your New Password", func() (string, error) {
		return utils.LoadRecoveryEmailTemplate(setupPassLink)
	})
}

func (ru *resendUseCase) SendContactNotifyMail(env *config.Env, contact *domain.Contact) (string, error) {
	return ru.sendEmail(env.NotifyEmail, "📬 SmartSRT - New Contact Form Submission", func() (string, error) {
		return utils.LoadContactNotifyTemplate(contact)
	})
}

func (ru *resendUseCase) SendDeleteAccountEmail(email, deleteAccountLink string) (string, error) {
	return ru.sendEmail(email, "⚠️ SmartSRT - Confirm Account Deletion", func() (string, error) {
		return utils.LoadDeleteAccountEmailTemplate(deleteAccountLink)
	})
}

func (ru *resendUseCase) SendSRTCreatedEmail(email, SRTLink string) (string, error) {
	return ru.sendEmail(email, "✅ SmartSRT - Your Subtitles Are Ready!", func() (string, error) {
		return utils.LoadSRTCreatedEmailTemplate(SRTLink)
	})
}
