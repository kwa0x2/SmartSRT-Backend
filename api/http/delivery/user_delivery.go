package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type UserDelivery struct {
	UserUseCase        domain.UserUseCase
	UserBaseRepository domain.BaseRepository[*domain.User]
}

func (ud *UserDelivery) GetProfileFromSession(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)
	userData.Password = ""

	ctx.JSON(http.StatusOK, userData)
}

func (ud *UserDelivery) CheckEmailExists(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	exists, err := ud.UserUseCase.IsEmailExists(email)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "email_existence_check", "email": email})
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if exists {
		ctx.Status(http.StatusFound)
	} else {
		ctx.Status(http.StatusOK) // available
	}
}

func (ud *UserDelivery) CheckPhoneExists(ctx *gin.Context) {
	phone := ctx.Param("phone")
	if phone == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	exists, err := ud.UserUseCase.IsPhoneExists(phone)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "phone_existence_check", "phone": phone})
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if exists {
		ctx.Status(http.StatusFound)
	} else {
		ctx.Status(http.StatusOK) // available
	}
}
