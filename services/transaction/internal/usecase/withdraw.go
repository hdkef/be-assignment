package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

// Withdraw implements usecase.TransactionUsecase.
func (t *TransactionUC) Withdraw(ctx context.Context, dto *entity.WithdrawTransactionDto) error {

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

	// update balance
	err = t.AccBalanceRepo.DecrementBalance(ctx, dto.AccID, dto.Amount, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// insert transaction logs
	trxLog := entity.TransactionLogs{
		ReffNum:   uuid.New(),
		CreatedAt: time.Now(),
		AccID:     dto.AccID,
		Amount:    dto.Amount,
		Status:    entity.ENUM_TRX_LOGS_STATUS_SUCCESS,
		Desc:      dto.Desc,
		ToAccID:   nil,
		EventType: entity.ENUM_TRX_LOGS_EVENT_TYPE_WITHDRAW,
	}

	err = t.TrxLogsRepo.CreateLogs(ctx, &trxLog, uow)
	if err != nil {
		logger.LogError(ctx, err)
		uow.Tx.Rollback()
		return err
	}

	// publish transaction created event

	trxCreatedEvent := et2.TransactionCreatedEventDto{
		ReffNum:   trxLog.ReffNum.String(),
		EventType: et2.ENUM_TRX_CREATED_EVENT_TYPE_WITHDRAW,
		Detail: et2.TransactionCreatedEventDtoDetail{
			AccID:     trxLog.AccID.String(),
			Amount:    trxLog.Amount,
			Status:    et2.ENUM_TRX_CREATED_STATUS(trxLog.Status),
			Desc:      trxLog.Desc,
			CreatedAt: trxLog.CreatedAt,
		},
	}

	err = t.Publisher.PublishCreateTransactionEvents(ctx, &trxCreatedEvent)
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
