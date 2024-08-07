package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/repository"
)

type TransactionLogsRepo struct {
	db *sql.DB
}

// CreateLogs implements repository.TransactionLogsRepository.
func (t *TransactionLogsRepo) CreateLogs(ctx context.Context, logs *entity.TransactionLogs, uow *repository.UnitOfWork) error {
	query := `
        INSERT INTO transactions.transactions_logs (
            reff_num, created_at, acc_id, amount, description, status, to_acc_id, event_type
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    `

	var err error
	if uow != nil && uow.Tx != nil {
		_, err = uow.Tx.ExecContext(ctx, query, logs.ReffNum, logs.CreatedAt, logs.AccID, logs.Amount, logs.Desc, string(logs.Status), logs.ToAccID, logs.EventType)
	} else {
		_, err = t.db.ExecContext(ctx, query, logs.ReffNum, logs.CreatedAt, logs.AccID, logs.Amount, logs.Desc, string(logs.Status), logs.ToAccID, logs.EventType)
	}

	if err != nil {
		return fmt.Errorf("failed to create transaction log: %w", err)
	}

	return nil
}

func NewTransactionLogsRepo(db *sql.DB) repository.TransactionLogsRepository {
	return &TransactionLogsRepo{
		db: db,
	}
}
