package delivery

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
)

type ContactDelivery struct {
	ContactUseCase domain.ContactUseCase
	ResendUseCase  domain.ResendUseCase
	Env            *config.Env
}

func (cd *ContactDelivery) Create(ctx *gin.Context) {
	var body domain.ContactCreateBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to bind JSON for contact form",
				slog.String("action", "validation_contact_form"),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	contact := &domain.Contact{
		FirstName: body.FirstName,
		LastName:  body.LastName,
		Email:     body.Email,
		Message:   body.Message,
	}

	if err := cd.ContactUseCase.Create(contact); err != nil {
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to save contact form to database",
				slog.String("action", "contact_form_database_save"),
				slog.String("email", body.Email),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	_, sentErr := cd.ResendUseCase.SendContactNotifyMail(cd.Env, contact)
	if sentErr != nil {
		if !utils.IsNormalBusinessError(sentErr) {
			slog.Error("Failed to send contact notification email",
				slog.String("action", "contact_notification_email_sending"),
				slog.String("email", body.Email),
				slog.String("error", sentErr.Error()))
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send new contact form email. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("Message sent successfully!"))
}
