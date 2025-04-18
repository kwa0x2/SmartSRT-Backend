package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
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

func SetSIDCookie(ctx *gin.Context, sessionID string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "sid",
		Value:    sessionID,
		MaxAge:   259200, // 3 day
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func DeleteCookie(ctx *gin.Context, name string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteLaxMode,
	})
}

func SetAuthTokenCookie(ctx *gin.Context, token, path string, maxAge int) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   false,
		Path:     path,
		SameSite: http.SameSiteLaxMode,
	})
}

func SetErrorCookie(ctx *gin.Context, value, path string) {
	http.SetCookie(ctx.Writer, &http.Cookie{
		Name:     "error",
		Value:    value,
		MaxAge:   15,
		HttpOnly: false,
		Secure:   false,
		Path:     path,
		SameSite: http.SameSiteLaxMode,
	})
}
