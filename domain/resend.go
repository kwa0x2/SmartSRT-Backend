package domain

type ResendRepository interface {
	SendEmail(to, subject, htmlContent string) (string, error)
}

type ResendUseCase interface {
	SendRecoveryEmail(email, recoveryLink string) (string, error)
}
