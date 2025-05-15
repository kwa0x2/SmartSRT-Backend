package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"time"
)

type customerUseCase struct {
	customerBaseRepository domain.BaseRepository[*domain.Customer]
}

func NewCustomerUseCase(customerBaseRepository domain.BaseRepository[*domain.Customer]) domain.CustomerUsaCase {
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
