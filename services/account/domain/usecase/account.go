package usecase

import (
	"context"

	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
)

type AccountUC interface {
	TransactionCreatedSend(ctx context.Context, dto et2.TransactionCreatedEventDto) error
	TransactionCreatedWithdraw(ctx context.Context, dto et2.TransactionCreatedEventDto) error
	GetHistory(ctx context.Context, dto *entity.GetHistoryDto) ([]*entity.History, error)
	CreateAccount(ctx context.Context, dto *entity.CreateAccountDto) error
}
