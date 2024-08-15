package repository

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

type QueueRepo struct {
	db *sql.DB
}

// GetPending implements repository.QueueRepository.
func (q *QueueRepo) GetPending(ctx context.Context, uow *repository.UnitOfWork) ([]entity.QueuedTrx, error) {
	panic("unimplemented")
}

// SetStatus implements repository.QueueRepository.
func (q *QueueRepo) SetStatus(ctx context.Context, id uuid.UUID, status string, uow *repository.UnitOfWork) error {
	panic("unimplemented")
}

// Create implements repository.QueueRepository.
func (*QueueRepo) Create(ctx context.Context, q *entity.QueuedTrx, uow *repository.UnitOfWork) error {
	panic("unimplemented")
}

func NewQueueRepo(db *sql.DB) repository.QueueRepository {
	return &QueueRepo{
		db: db,
	}
}
