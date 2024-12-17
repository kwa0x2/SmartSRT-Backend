package usecase

import (
	"context"
	"errors"
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

	TTL := time.Now().UTC().Add(2 * time.Minute).Unix()

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
		return errors.New("session has expired")
	}

	newTTL := time.Now().UTC().Add(2 * time.Minute).Unix()

	if err = su.sessionRepository.UpdateSessionTTL(ctx, sessionID, int(newTTL)); err != nil {
		return err
	}

	return nil
}
