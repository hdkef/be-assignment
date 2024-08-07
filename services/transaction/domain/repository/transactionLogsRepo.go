package repository

import (
	"context"

	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type TransactionLogsRepository interface {
	CreateLogs(ctx context.Context, logs *entity.TransactionLogs, uow *UnitOfWork) error
}
