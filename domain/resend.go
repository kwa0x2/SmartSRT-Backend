package domain

type ResendRepository interface {
	SendEmail(to, subject, htmlContent string) (string, error)
}

type ResendUseCase interface {
	SendSetupPasswordEmail(email, setupPassLink string) (string, error)
}
