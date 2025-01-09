package utils

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"time"
)

var secretKey []byte

func GenerateJWT(jwtClaims jwt.MapClaims, env *bootstrap.Env) (string, error) {
	secretKey = []byte(env.JWTSecret)

	expirationTime := time.Now().Add(1 * time.Hour).Unix() // 1 hour
	jwtClaims["exp"] = expirationTime

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
