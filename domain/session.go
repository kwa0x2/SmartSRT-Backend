package domain

import (
	"context"
)

const (
	TableName = "sessions"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, sessionID string, TTL int) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	UpdateSessionTTL(ctx context.Context, sessionID string, newTTL int) error
}

type SessionUseCase interface {
	CreateSession() (string, error)
	ValidateSession(sessionID string) error
}

type Session struct {
	SessionID string
	TTL       int
}
