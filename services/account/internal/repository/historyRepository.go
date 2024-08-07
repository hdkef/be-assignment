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

// GetHistory implements repository.HistoryRepository.
func (h *HistoryRepo) GetHistory(ctx context.Context, f *entity.GetHistoryFilter, opt *entity.GetHistoryOptions) ([]*entity.History, error) {
	query := `
	SELECT id, reff_num, created_at, acc_id, trx_type, amount, status, description, acc_id_2
	FROM accounts.histories
	WHERE acc_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
`

	offset := (opt.Page - 1) * opt.Limit

	rows, err := h.db.QueryContext(ctx, query, f.AccID, opt.Limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query history records: %w", err)
	}
	defer rows.Close()

	var histories []*entity.History
	for rows.Next() {
		var h entity.History
		var accID2 sql.NullString
		err := rows.Scan(&h.Id, &h.ReffNum, &h.CreatedAt, &h.AccID, &h.TrxType, &h.Amount, &h.Status, &h.Desc, &accID2)
		if err != nil {
			return nil, fmt.Errorf("failed to scan history record: %w", err)
		}
		if accID2.Valid {
			h.AccID2, err = uuid.Parse(accID2.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse AccID2: %w", err)
			}
		}
		histories = append(histories, &h)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating history records: %w", err)
	}

	return histories, nil
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
