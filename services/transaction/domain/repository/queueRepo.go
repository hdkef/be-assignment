package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type QueueRepository interface {
	GetPending(ctx context.Context, uow *UnitOfWork) ([]entity.QueuedTrx, error)
	Create(ctx context.Context, q *entity.QueuedTrx, uow *UnitOfWork) error
	SetStatus(ctx context.Context, id uuid.UUID, status string, uow *UnitOfWork) error
}
