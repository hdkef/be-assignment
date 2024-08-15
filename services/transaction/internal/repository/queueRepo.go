package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

type QueueRepo struct {
	db *sql.DB
}

// GetPending implements repository.QueueRepository.
func (q *QueueRepo) GetPending(ctx context.Context, uow *repository.UnitOfWork) ([]entity.QueuedTrx, error) {
	query := `
		SELECT id, created_at, status, result, schedule_trx_id
		FROM transactions.queued_trx
		WHERE status = $1
	`

	var rows *sql.Rows
	var err error
	if uow != nil {
		rows, err = uow.Tx.QueryContext(ctx, query, entity.ENUM_QUEUED_TRX_STATUS_PENDING)
	} else {
		rows, err = q.db.QueryContext(ctx, query, entity.ENUM_QUEUED_TRX_STATUS_PENDING)
	}

	if err != nil {
		return nil, fmt.Errorf("error querying pending transactions: %w", err)
	}
	defer rows.Close()

	var transactions []entity.QueuedTrx
	for rows.Next() {
		var trx entity.QueuedTrx
		err := rows.Scan(
			&trx.ID, &trx.CreatedAt, &trx.Status, &trx.Result, &trx.ScheduleTrxID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		transactions = append(transactions, trx)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return transactions, nil
}

// SetStatus implements repository.QueueRepository.
func (q *QueueRepo) SetStatus(ctx context.Context, id uuid.UUID, status entity.ENUM_QUEUED_TRX_STATUS, uow *repository.UnitOfWork) error {
	query := `
		UPDATE transactions.queued_trx
		SET status = $1
		WHERE id = $2
	`

	var result sql.Result
	var err error
	if uow != nil {
		result, err = uow.Tx.ExecContext(ctx, query, status, id)
	} else {
		result, err = q.db.ExecContext(ctx, query, status, id)
	}

	if err != nil {
		return fmt.Errorf("error updating queued transaction status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no queued transaction found with id %s", id)
	}

	return nil
}

// SetResult implements repository.QueueRepository.
func (q *QueueRepo) SetResult(ctx context.Context, id uuid.UUID, res entity.ENUM_QUEUED_TRX_RESULT, uow *repository.UnitOfWork) error {
	query := `
		UPDATE transactions.queued_trx
		SET result = $1
		WHERE id = $2
	`

	var result sql.Result
	var err error
	if uow != nil {
		result, err = uow.Tx.ExecContext(ctx, query, res, id)
	} else {
		result, err = q.db.ExecContext(ctx, query, res, id)
	}

	if err != nil {
		return fmt.Errorf("error updating queued transaction status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no queued transaction found with id %s", id)
	}

	return nil
}

// Create implements repository.QueueRepository.
func (q *QueueRepo) Create(ctx context.Context, trx *entity.QueuedTrx, uow *repository.UnitOfWork) error {
	query := `
		INSERT INTO transactions.queued_trx (id, created, status, result, schedule_trx_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	var result sql.Result
	var err error
	if uow != nil {
		result, err = uow.Tx.ExecContext(ctx, query, trx.ID, trx.CreatedAt, trx.Status, trx.Result, trx.ScheduleTrxID)
	} else {
		result, err = q.db.ExecContext(ctx, query, trx.ID, trx.CreatedAt, trx.Status, trx.Result, trx.ScheduleTrxID)
	}

	if err != nil {
		return fmt.Errorf("error creating queued transaction: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no queued transaction was created")
	}

	return nil
}

func NewQueueRepo(db *sql.DB) repository.QueueRepository {
	return &QueueRepo{
		db: db,
	}
}
