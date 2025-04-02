package domain

import "github.com/kwa0x2/AutoSRT-Backend/bootstrap"

type ResendRepository interface {
	SendEmail(to, subject, htmlContent string) (string, error)
}

type ResendUseCase interface {
	SendSetupPasswordEmail(email, setupPassLink string) (string, error)
	SendContactNotifyMail(env *bootstrap.Env, contact *Contact) (string, error)
	SendDeleteAccountEmail(email, deleteAccountLink string) (string, error)
}
