package http

import "github.com/hdkef/be-assignment/services/transaction/domain/usecase"

type HttpHandler struct {
	TransactionUC usecase.TransactionUsecase
}
