package utils

import "errors"

var ErrSessionExpired = errors.New("session is expired")
var ErrSessionNotFound = errors.New("session not found in dynamodb")
var ErrLimitReached = errors.New("monthly usage limit reached")
