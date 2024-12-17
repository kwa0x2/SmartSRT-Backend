package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
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
