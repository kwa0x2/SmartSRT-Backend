package repository

import (
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/resend/resend-go/v2"
)

type ResendRepository struct {
	resendClient *resend.Client
}

func NewResendRepository(resendClient *resend.Client) domain.ResendRepository {
	return &ResendRepository{
		resendClient: resendClient,
	}
}

func (rr *ResendRepository) SendEmail(to, subject, htmlContent string) (string, error) {
	params := &resend.SendEmailRequest{
		From:    "SmartSRT <no-reply@alperkarakoyun.com>",
		To:      []string{to},
		Html:    htmlContent,
		Subject: subject,
	}

	sent, err := rr.resendClient.Emails.Send(params)
	if err != nil {
		return "", err
	}
	return sent.Id, nil
}
