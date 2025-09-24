package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
)

func LocaleMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		locale := utils.GetLocale(ctx)
		ctx.Set("locale", locale)
		ctx.Next()
	}
}
