package utils

import "errors"

var ErrTTLMissing = errors.New("TTL field is missing or invalid")

var ErrSessionExpired = errors.New("session is expired")
