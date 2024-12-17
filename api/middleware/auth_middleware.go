package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"net/http"
	"time"
)

func SessionMiddleware(sessionUseCase domain.SessionUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sessionID, err := ctx.Cookie("sid")
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", "Session ID cookie is required"))
			ctx.Abort()
			return
		}

		err = sessionUseCase.ValidateSession(sessionID)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, utils.NewErrorResponse("Unauthorized", err.Error()))
			ctx.Abort()
			return
		}

		newExpiryTime := time.Now().UTC().Add(2 * time.Minute)

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "sid",
			Value:    sessionID,
			Expires:  newExpiryTime,
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		})

		ctx.Next()
	}
}
