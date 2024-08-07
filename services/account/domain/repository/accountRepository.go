package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type AccountRepository interface {
	CreateAccount(ctx context.Context, acc *entity.Account, uow *UnitOfWork) error
	IncrementBalance(ctx context.Context, id uuid.UUID, incrAmount float64, uow *UnitOfWork) error
	DecrementBalance(ctx context.Context, id uuid.UUID, decrAmount float64, uow *UnitOfWork) error
}
