package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	TableName = "sessions"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, sessionID string) (*Session, error)
	UpdateSessionTTL(ctx context.Context, sessionID string, newTTL int) error
	DeleteSession(ctx context.Context, sessionID string) error
}

type SessionUseCase interface {
	CreateSession(userID bson.ObjectID) (string, error)
	ValidateSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
}

type Session struct {
	SessionID string `dynamodbav:"session_id"`
	UserID    string `dynamodbav:"user_id"`
	Role      string `dynamodbav:"role"`
	TTL       int    `dynamodbav:"ttl"`
}
