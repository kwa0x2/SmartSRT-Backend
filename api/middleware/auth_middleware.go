package middleware

import (
	"net/http"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

func SessionMiddleware(sessionUseCase domain.SessionUseCase, userBaseRepository domain.BaseRepository[*domain.User]) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("sid")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Session ID cookie is required"))
			ctx.Abort()
			return
		}

		session, err := sessionUseCase.ValidateSession(sessionID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse(err.Error()))
			ctx.Abort()
			return
		}

		userID, err := bson.ObjectIDFromHex(session.UserID)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
			ctx.Abort()
			return
		}

		filter := bson.D{{Key: "_id", Value: userID}}
		result, err := userBaseRepository.FindOne(nil, filter)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
			ctx.Abort()
			return
		}

		utils.SetSIDCookie(ctx, sessionID)
		result.ID = userID
		ctx.Set("user", result)

		ctx.Next()
	}
}

func JWTMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Authorization header is missing. Please try again later or contact support."))
			ctx.Abort()
			return
		}

		claims, err := utils.GetClaims(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized2. Please try again later or contact support."))
			ctx.Abort()
			return
		}

		ctx.Set("claims", claims)
		ctx.Next()
	}
}
