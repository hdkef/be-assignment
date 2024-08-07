package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type AccountBalanceRepository interface {
	CreateAccountBalance(ctx context.Context, accBalance *entity.AccountBalance, uow *UnitOfWork) error
	IncrementBalance(ctx context.Context, id uuid.UUID, incrAmount float64, uow *UnitOfWork) error
	DecrementBalance(ctx context.Context, id uuid.UUID, decrAmount float64, uow *UnitOfWork) error
}
