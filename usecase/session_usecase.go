package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
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

func (su *sessionUseCase) CreateSession() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sessionID, err := utils.GenerateSessionID()
	if err != nil {
		return "", err
	}

	TTL := time.Now().UTC().Add(24 * time.Hour).Unix()

	err = su.sessionRepository.CreateSession(ctx, sessionID, int(TTL))
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (su *sessionUseCase) ValidateSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := su.sessionRepository.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	currentTimeUnix := time.Now().UTC().Unix()

	if currentTimeUnix > int64(session.TTL) {
		return utils.ErrSessionExpired
	}

	newTTL := time.Now().UTC().Add(24 * time.Hour).Unix()

	if err = su.sessionRepository.UpdateSessionTTL(ctx, sessionID, int(newTTL)); err != nil {
		return err
	}

	return nil
}

func (su *sessionUseCase) DeleteSession(sessionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := su.sessionRepository.DeleteSession(ctx, sessionID); err != nil {
		return err
	}

	return nil
}
