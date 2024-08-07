package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

// Send implements usecase.TransactionUsecase.
func (t *TransactionUC) Send(ctx context.Context, dto *entity.SendTransactionDto) error {

	err := dto.Validate()
	if err != nil {
		return err
	}

	// begin trx
	uow, err := t.UoW.NewUnitOfWork()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	// update balance
	err = t.AccBalanceRepo.DecrementBalance(ctx, dto.AccID, dto.Amount, uow)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	err = t.AccBalanceRepo.IncrementBalance(ctx, dto.ToAccID, dto.Amount, uow)
	if err != nil {
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
		ToAccID:   &dto.ToAccID,
		EventType: entity.ENUM_TRX_LOGS_EVENT_TYPE_SEND,
	}

	err = t.TrxLogsRepo.CreateLogs(ctx, &trxLog, uow)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	// publish transaction created event
	toAccId := trxLog.ToAccID.String()

	trxCreatedEvent := et2.TransactionCreatedEventDto{
		ReffNum:   trxLog.ReffNum.String(),
		EventType: et2.ENUM_TRX_CREATED_EVENT_TYPE_SEND,
		Detail: et2.TransactionCreatedEventDtoDetail{
			AccID:     trxLog.AccID.String(),
			Amount:    trxLog.Amount,
			Status:    et2.ENUM_TRX_CREATED_STATUS(trxLog.Status),
			Desc:      trxLog.Desc,
			CreatedAt: trxLog.CreatedAt,
			ToAccID:   &toAccId,
		},
	}

	err = t.Publisher.PublishCreateTransactionEvents(ctx, &trxCreatedEvent)
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	err = uow.Tx.Commit()
	if err != nil {
		uow.Tx.Rollback()
		return err
	}

	return nil
}
