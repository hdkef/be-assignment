package service

import (
	"context"

	et2 "github.com/hdkef/be-assignment/pkg/domain/entity"
)

type Publisher interface {
	PublishCreateTransactionEvents(ctx context.Context, dto *et2.TransactionCreatedEventDto) error
}
