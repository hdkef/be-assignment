package usecase

import (
	"context"
	"fmt"

	"github.com/hdkef/be-assignment/pkg/logger"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

// ProcessQueue implements usecase.TransactionUsecase.
func (t *TransactionUC) ProcessQueue(ctx context.Context) error {

	// get queue status pending
	queues, err := t.QueueRepo.GetPending(ctx, nil)
	if err != nil {
		logger.LogError(ctx, err)
		return err
	}

	for _, q := range queues {

		// begin trx
		uow, err := t.UoW.NewUnitOfWork()
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			continue
		}

		// handle autodebet
		item, err := t.ScheduleRepo.Find(ctx, q.ScheduleTrxID, uow)
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			continue
		}

		var isSuccess bool = true

		switch item.Type {
		case entity.ENUM_SCHEDULE_TYPE_SEND:
			// handle send autodebet
			dto := entity.SendTransactionDto{
				AccID:   item.AccID,
				Amount:  item.Amount,
				Desc:    fmt.Sprintf("from autodebetId %s", item.ID.String()),
				ToAccID: item.ToAccID,
			}

			err = t.Send(ctx, &dto)
			if err != nil {
				logger.LogError(ctx, err)
				isSuccess = false
			}
		}

		result := entity.ENUM_QUEUED_TRX_RESULT_SUCCESS

		if !isSuccess {
			result = entity.ENUM_QUEUED_TRX_RESULT_FAILED
		}

		// update status
		err = t.QueueRepo.SetStatus(ctx, q.ID, entity.ENUM_QUEUED_TRX_STATUS_EXECUTED, uow)
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			continue
		}

		// update result
		err = t.QueueRepo.SetResult(ctx, q.ID, result, uow)
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			continue
		}

		err = uow.Tx.Commit()
		if err != nil {
			logger.LogError(ctx, err)
			uow.Tx.Rollback()
			continue
		}
	}

	return nil
}
