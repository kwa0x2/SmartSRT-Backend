package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"time"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(userRepository domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uu *userUseCase) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	if err := user.Validate(); err != nil {
		return err
	}
	return uu.userRepository.Create(ctx, user)
}
