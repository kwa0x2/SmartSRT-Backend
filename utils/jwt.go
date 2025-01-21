package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
)

var secretKey []byte

func GenerateJWT(jwtClaims jwt.MapClaims, env *bootstrap.Env, expUnixTime int64) (string, error) {
	secretKey = []byte(env.JWTSecret)

	jwtClaims["exp"] = expUnixTime

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyJWT(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid token")
	}

	return nil
}

func GetClaims(tokenString string) (jwt.MapClaims, error) {
	err := VerifyJWT(tokenString)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{}
	_, parseErr := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secretKey, nil
	})

	if parseErr != nil {
		return nil, parseErr
	}

	return claims, nil
}
