package domain

import (
	"context"
	"github.com/kwa0x2/SmartSRT-Backend/domain/types"
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
	CreateSessionAndUpdateLastLogin(userID bson.ObjectID, plan types.PlanType, email string) (string, error)
	ValidateSession(sessionID string) (*Session, error)
	DeleteSession(sessionID string) error
}

type Session struct {
	SessionID string         `dynamodbav:"session_id"`
	UserID    string         `dynamodbav:"user_id"`
	Plan      types.PlanType `dynamodbav:"plan"`
	TTL       int            `dynamodbav:"ttl"`
}
