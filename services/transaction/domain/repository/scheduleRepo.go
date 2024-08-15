package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type ScheduleRepository interface {
	Create(ctx context.Context, sched *entity.ScheduleTrx, uow *UnitOfWork) error
	GetUnprocessedDaily(ctx context.Context, today time.Time, uow *UnitOfWork) ([]*entity.ScheduleTrx, error)
	UpdateChecked(ctx context.Context, ids []uuid.UUID, today time.Time, uow *UnitOfWork) error
	Find(ctx context.Context, id uuid.UUID, uow *UnitOfWork) (*entity.ScheduleTrx, error)
}
