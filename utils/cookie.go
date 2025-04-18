package utils

import (
	"github.com/gin-gonic/gin"
)

func GetLocale(ctx *gin.Context) string {
	locale, err := ctx.Cookie("NEXT_LOCALE")
	if err != nil {
		return "en"
	}

	if locale != "tr" && locale != "en" {
		return "en"
	}

	return locale
}
