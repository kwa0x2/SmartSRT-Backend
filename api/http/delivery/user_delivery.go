package delivery

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
)

type UserDelivery struct {
	UserUseCase domain.UserUseCase
}

func (ud *UserDelivery) GetProfileFromSession(ctx *gin.Context) {
	sessionUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	userIDStr, ok := sessionUserID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid session data. Please log in again. If the issue persists, contact support."))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred2. Please try again later or contact support."))
		return
	}

	user, err := ud.UserUseCase.FindOneByID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	user.Password = ""

	ctx.JSON(http.StatusOK, user)
}

func (ud *UserDelivery) CheckEmailExists(ctx *gin.Context) {
	email := ctx.Param("email")

	if email == "" {
		ctx.Status(http.StatusInternalServerError)
		return
	}

	exists, err := ud.UserUseCase.IsEmailExists(email)
	if err != nil {
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
		ctx.Status(http.StatusInternalServerError)
		return
	}

	if exists {
		ctx.Status(http.StatusFound)
	} else {
		ctx.Status(http.StatusOK) // available
	}
}
