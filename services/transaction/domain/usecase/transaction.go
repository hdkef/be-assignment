package usecase

import (
	"context"

	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

type TransactionUsecase interface {
	Send(ctx context.Context, dto *entity.SendTransactionDto) error
	Withdraw(ctx context.Context, dto *entity.WithdrawTransactionDto) error
}
