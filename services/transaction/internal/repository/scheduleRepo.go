package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
	"github.com/lib/pq"
)

type ScheduleRepo struct {
	db *sql.DB
}

// Find implements repository.ScheduleRepository.
func (s *ScheduleRepo) Find(ctx context.Context, id uuid.UUID, uow *repository.UnitOfWork) (*entity.ScheduleTrx, error) {
	query := `
		SELECT id, created_at, status, acc_id, type, schedule, to_acc_id, amount, has_checked, last_checked
		FROM transactions.scheduled_trx
		WHERE id = $1
	`

	var row *sql.Row
	if uow != nil {
		row = uow.Tx.QueryRowContext(ctx, query, id)
	} else {
		row = s.db.QueryRowContext(ctx, query, id)
	}

	var trx entity.ScheduleTrx
	err := row.Scan(
		&trx.ID, &trx.CreatedAt, &trx.Status, &trx.AccID, &trx.Type, &trx.Schedule,
		&trx.ToAccID, &trx.Amount, &trx.HasChecked, &trx.LastChecked,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("scheduled transaction not found")
		}
		return nil, fmt.Errorf("error finding scheduled transaction: %w", err)
	}

	return &trx, nil
}

// GetUnprocessed implements repository.ScheduleRepository.
// find last checked is < today and status active
func (s *ScheduleRepo) GetUnprocessed(ctx context.Context, today time.Time, uow *repository.UnitOfWork) ([]*entity.ScheduleTrx, error) {
	query := `
		SELECT id, created_at, status, acc_id, type, schedule, to_acc_id, amount, has_checked, last_checked
		FROM transactions.scheduled_trx
		WHERE status = $1 AND (last_checked IS NULL OR last_checked < $2)
	`

	var rows *sql.Rows
	var err error
	if uow != nil {
		rows, err = uow.Tx.QueryContext(ctx, query, entity.ENUM_SCHEDULE_TRX_STATUS_ACTIVE, today)
	} else {
		rows, err = s.db.QueryContext(ctx, query, entity.ENUM_SCHEDULE_TRX_STATUS_ACTIVE, today)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying unprocessed transactions: %w", err)
	}
	defer rows.Close()

	var transactions []*entity.ScheduleTrx
	for rows.Next() {
		var trx entity.ScheduleTrx
		err := rows.Scan(
			&trx.ID, &trx.CreatedAt, &trx.Status, &trx.AccID, &trx.Type, &trx.Schedule,
			&trx.ToAccID, &trx.Amount, &trx.HasChecked, &trx.LastChecked,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		transactions = append(transactions, &trx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

// UpdateChecked implements repository.ScheduleRepository.
// update last checked today
func (s *ScheduleRepo) UpdateChecked(ctx context.Context, ids []uuid.UUID, today time.Time, uow *repository.UnitOfWork) error {
	query := `
		UPDATE transactions.scheduled_trx
		SET last_checked = $1, has_checked = true
		WHERE id = ANY($2)
	`

	var err error
	if uow != nil {
		_, err = uow.Tx.ExecContext(ctx, query, today, pq.Array(ids))
	} else {
		_, err = s.db.ExecContext(ctx, query, today, pq.Array(ids))
	}

	if err != nil {
		return fmt.Errorf("error updating checked status: %w", err)
	}

	return nil
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
