package usecase

import (
	"context"

	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type TransactionUsecase interface {
	Send(ctx context.Context, dto *entity.SendTransactionDto) error
	Withdraw(ctx context.Context, dto *entity.WithdrawTransactionDto) error
	CreateAutodebet(ctx context.Context, dto *entity.CreateAutodebetDto) error
	ProcessAutodebetDaily(ctx context.Context) error
	ProcessQueue(ctx context.Context) error
}
