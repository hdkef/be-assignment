package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/repository"
)

type HistoryRepo struct {
	db *sql.DB
}

// CreateHistory implements repository.HistoryRepository.
func (h *HistoryRepo) CreateHistory(ctx context.Context, history *entity.History, uow *repository.UnitOfWork) error {
	query := `
        INSERT INTO accounts.histories (
            id, reff_num, created_at, acc_id, trx_type, amount, status, description, acc_id_2
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    `
	var accid2 *uuid.UUID

	if uuid.Nil != history.AccID2 {
		accid2 = &history.AccID2
	} else {
		accid2 = nil
	}

	var err error
	if uow != nil && uow.Tx != nil {
		_, err = uow.Tx.ExecContext(ctx, query, history.Id, history.ReffNum, history.CreatedAt, history.AccID, history.TrxType, history.Amount, history.Status, history.Desc, accid2)
	} else {
		_, err = h.db.ExecContext(ctx, query, history.Id, history.ReffNum, history.CreatedAt, history.AccID, history.TrxType, history.Amount, history.Status, history.Desc, accid2)
	}

	if err != nil {
		return fmt.Errorf("failed to create history record: %w", err)
	}

	return nil
}

func NewHistoryRepo(db *sql.DB) repository.HistoryRepository {
	return &HistoryRepo{
		db: db,
	}
}
