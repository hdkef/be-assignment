package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

// CreateAutodebet implements usecase.TransactionUsecase.
func (t *TransactionUC) CreateAutodebet(ctx context.Context, dto *entity.CreateAutodebetDto) error {
	err := dto.Validate()
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	// begin trx
	uow, err := t.UoW.NewUnitOfWork()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// insert autodebet
	trxLog := entity.ScheduleTrx{
		ID:          uuid.New(),
		CreatedAt:   time.Now(),
		Status:      entity.ENUM_SCHEDULE_TRX_STATUS_ACTIVE,
		Type:        entity.ENUM_SCHEDULE_TYPE(dto.Type),
		Schedule:    dto.Schedule,
		ToAccID:     dto.ToAccID,
		AccID:       dto.AccID,
		Amount:      dto.Amount,
		HasChecked:  false,
		LastChecked: time.Now(),
	}

	err = t.ScheduleRepo.Create(ctx, &trxLog, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	err = uow.Tx.Commit()
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	return nil
}
