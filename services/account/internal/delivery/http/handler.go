package http

import "github.com/hdkef/be-assignment/services/account/domain/usecase"

type HttpHandler struct {
	UserUsecase usecase.UserUsecase
	AccUsecase  usecase.AccountUC
}
