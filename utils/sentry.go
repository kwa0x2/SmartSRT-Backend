package utils

import (
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// IsNormalBusinessError checks if an error is a normal business logic error that shouldn't be sent to Sentry
func IsNormalBusinessError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// MongoDB normal errors
	if errors.Is(err, mongo.ErrNoDocuments) {
		return true
	}

	// Session related normal errors
	if errors.Is(err, ErrSessionExpired) || errors.Is(err, ErrSessionNotFound) {
		return true
	}

	// Common normal error patterns
	normalErrorPatterns := []string{
		"session not found",
		"user not found",
		"invalid credentials",
		"unauthorized",
		"validation failed",
		"duplicate key",
		"record not found",
	}

	for _, pattern := range normalErrorPatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	return false
}
