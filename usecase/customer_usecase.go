package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type customerUseCase struct {
	customerBaseRepository domain.BaseRepository[*domain.Customer]
}

func NewCustomerUseCase(customerBaseRepository domain.BaseRepository[*domain.Customer]) domain.CustomerUseCase {
	return &customerUseCase{customerBaseRepository: customerBaseRepository}
}

func (cu *customerUseCase) Create(customer domain.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	return cu.customerBaseRepository.Create(ctx, &customer)
}

func (cu *customerUseCase) FindByEmail(email string) (*domain.Customer, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}

	return cu.customerBaseRepository.FindOne(ctx, filter)
}
