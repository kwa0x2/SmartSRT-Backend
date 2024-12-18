package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"net/http"
)

func SessionMiddleware(sessionUseCase domain.SessionUseCase) gin.HandlerFunc {
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

		http.SetCookie(ctx.Writer, &http.Cookie{
			Name:     "sid",
			Value:    sessionID,
			MaxAge:   86400, // 24 hours
			HttpOnly: true,
			Secure:   false,
			Path:     "/",
		})

		ctx.Set("user_id", session.UserID)
		ctx.Set("role", session.Role)
		ctx.Next()
	}
}
