package utils

import (
	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-gonic/gin"
)

func CaptureError(err error, ctx *gin.Context, additionalData ...map[string]interface{}) {
	if hub := sentrygin.GetHubFromContext(ctx); hub != nil {
		hub.WithScope(func(scope *sentry.Scope) {
			if user, exists := ctx.Get("user"); exists {
				if userData, ok := user.(map[string]interface{}); ok {
					scope.SetUser(sentry.User{
						ID:    getStringFromMap(userData, "id"),
						Email: getStringFromMap(userData, "email"),
					})
				}
			}

			scope.SetTag("endpoint", ctx.FullPath())
			scope.SetTag("method", ctx.Request.Method)
			scope.SetTag("status_code", string(rune(ctx.Writer.Status())))

			for _, data := range additionalData {
				for key, value := range data {
					scope.SetExtra(key, value)
				}
			}

			hub.CaptureException(err)
		})
	}
}

func getStringFromMap(m map[string]interface{}, key string) string {
	if val, exists := m[key]; exists {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}
