package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/config"
)

func GenerateSessionID() (string, error) {
	bytes := make([]byte, 32)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate session ID: %w", err)
	}

	sessionID := base64.URLEncoding.EncodeToString(bytes)

	return sessionID, nil
}

func SetSIDCookie(ctx *gin.Context, sessionID string, env *config.Env) {
	isSecure := env.AppEnv == "production"
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		MaxAge:   259200, // 3 day
		HttpOnly: true,
		Secure:   isSecure,
		Path:     "/",
		Domain:   ".smartsrt.com",
		SameSite: http.SameSiteLaxMode,
	})
}

func DeleteCookie(ctx *gin.Context, name string, path *string, env *config.Env) {
	cookiePath := "/"
	if path != nil {
		cookiePath = *path
	}

	isSecure := env.AppEnv == "production"
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		Path:     cookiePath,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAuthTokenCookie(ctx *gin.Context, token, path string, maxAge int, env *config.Env) {
	isSecure := env.AppEnv == "production"
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isSecure,
		Path:     path,
		Domain:   ".smartsrt.com",
		SameSite: http.SameSiteLaxMode,
	})
}

func SetErrorCookie(ctx *gin.Context, value string, env *config.Env) {
	isSecure := env.AppEnv == "production"
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "error",
		Value:    value,
		MaxAge:   15,
		HttpOnly: false,
		Secure:   isSecure,
		Path:     "/",
		Domain:   ".smartsrt.com",
		SameSite: http.SameSiteLaxMode,
	})
}
