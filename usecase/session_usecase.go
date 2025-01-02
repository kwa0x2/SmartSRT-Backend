package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type sessionUseCase struct {
	sessionRepository domain.SessionRepository
}

func NewSessionUseCase(sessionRepository domain.SessionRepository) domain.SessionUseCase {
	return &sessionUseCase{
		sessionRepository: sessionRepository,
	}
}

func (su *sessionUseCase) CreateSession(userID bson.ObjectID) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return "", err
	}

	TTL := time.Now().UTC().Add(24 * time.Hour).Unix()

	session := domain.Session{
		SessionID: sessionID,
		UserID:    userID.Hex(),
		TTL:       int(TTL),
	}

	err = su.sessionRepository.CreateSession(ctx, session)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (su *sessionUseCase) ValidateSession(sessionID string) (*domain.Session, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := su.sessionRepository.GetSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	currentTimeUnix := time.Now().UTC().Unix()

	if currentTimeUnix > int64(session.TTL) {
		return nil, utils.ErrSessionExpired
	}

	newTTL := time.Now().UTC().Add(24 * time.Hour).Unix()

	if err = su.sessionRepository.UpdateSessionTTL(ctx, sessionID, int(newTTL)); err != nil {
		return nil, err
	}

	return session, nil
}

func (su *sessionUseCase) DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := su.sessionRepository.DeleteSession(ctx, sessionID); err != nil {
		return err
	}

	return nil
}
