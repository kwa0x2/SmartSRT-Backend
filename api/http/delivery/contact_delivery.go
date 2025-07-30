package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type ContactDelivery struct {
	ContactUseCase domain.ContactUseCase
	ResendUseCase  domain.ResendUseCase
	Env            *config.Env
}

func (cd *ContactDelivery) Create(ctx *gin.Context) {
	var body domain.ContactCreateBody

	if err := ctx.ShouldBindJSON(&body); err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "json_binding_contact_create"})
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
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "contact_creation"})
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	_, sentErr := cd.ResendUseCase.SendContactNotifyMail(cd.Env, contact)
	if sentErr != nil {
		utils.HandleErrorWithSentry(ctx, sentErr, map[string]interface{}{"action": "contact_notification_email"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to send new contact form email. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, utils.NewMessageResponse("Message sent successfully!"))
}
