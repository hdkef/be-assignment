package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type QueueRepository interface {
	GetPending(ctx context.Context, uow *UnitOfWork) ([]entity.QueuedTrx, error)
	Create(ctx context.Context, q *entity.QueuedTrx, uow *UnitOfWork) error
	SetStatus(ctx context.Context, id uuid.UUID, status entity.ENUM_QUEUED_TRX_STATUS, uow *UnitOfWork) error
	SetResult(ctx context.Context, id uuid.UUID, result entity.ENUM_QUEUED_TRX_RESULT, uow *UnitOfWork) error
}
