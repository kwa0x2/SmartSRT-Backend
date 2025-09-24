package delivery

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
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

	usage, exists := ctx.Get("usage")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)
	userData.Password = ""

	usageData := usage.(*domain.Usage)

	response := gin.H{
		"user":        userData,
		"usage_limit": usageData.UsageLimit,
	}

	ctx.JSON(http.StatusOK, response)
	
}

func (ud *UserDelivery) CheckEmailExists(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	exists, err := ud.UserUseCase.IsEmailExists(email)
	if err != nil {
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to check email existence",
				slog.String("action", "email_existence_check"),
				slog.String("email", email),
				slog.String("error", err.Error()))
		}
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
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to check phone existence",
				slog.String("action", "phone_existence_check"),
				slog.String("phone", phone),
				slog.String("error", err.Error()))
		}
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if exists {
		ctx.Status(http.StatusFound)
	} else {
		ctx.Status(http.StatusOK) // available
	}
}
