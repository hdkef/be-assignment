package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

type ScheduleRepo struct {
	db *sql.DB
}

// Find implements repository.ScheduleRepository.
func (s *ScheduleRepo) Find(ctx context.Context, id uuid.UUID, uow *repository.UnitOfWork) (*entity.ScheduleTrx, error) {
	panic("unimplemented")
}

// GetUnprocessed implements repository.ScheduleRepository.
func (s *ScheduleRepo) GetUnprocessed(ctx context.Context, today time.Time, uow *repository.UnitOfWork) ([]*entity.ScheduleTrx, error) {
	panic("unimplemented")
}

// UpdateChecked implements repository.ScheduleRepository.
func (s *ScheduleRepo) UpdateChecked(ctx context.Context, ids []uuid.UUID, today time.Time, uow *repository.UnitOfWork) error {
	panic("unimplemented")
}

// Create implements repository.ScheduleRepository.
func (s *ScheduleRepo) Create(ctx context.Context, sched *entity.ScheduleTrx, uow *repository.UnitOfWork) error {
	query := `
		INSERT INTO transactions.accounts_balance (id, created, status, acc_id, type, amount, schedule, to_acc_id, has_checked, last_checked)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	var err error
	if uow != nil && uow.Tx != nil {
		_, err = uow.Tx.ExecContext(ctx, query, sched.ID, sched.CreatedAt, sched.Status, sched.AccID, sched.Type, sched.Amount, sched.Schedule, sched.ToAccID, sched.HasChecked, sched.LastChecked)
	} else {
		_, err = s.db.ExecContext(ctx, query, sched.ID, sched.CreatedAt, sched.Status, sched.AccID, sched.Type, sched.Amount, sched.Schedule, sched.ToAccID, sched.HasChecked, sched.LastChecked)
	}

	if err != nil {
		return fmt.Errorf("failed to create account balance: %w", err)
	}

	return nil
}

func NewScheduleRepo(db *sql.DB) repository.ScheduleRepository {
	return &ScheduleRepo{
		db: db,
	}
}
