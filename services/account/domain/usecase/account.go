package usecase

import (
	"context"

	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
)

type AccountUC interface {
	TransactionCreatedSend(ctx context.Context, dto et2.TransactionCreatedEventDto) error
	TransactionCreatedWithdraw(ctx context.Context, dto et2.TransactionCreatedEventDto) error
}
