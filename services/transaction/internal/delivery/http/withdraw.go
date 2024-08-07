package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/delivery"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

func (t *HttpHandler) Withdraw(c *gin.Context) {

	req := struct {
		AccID  string  `json:"accId"`
		Desc   string  `json:"desc"`
		Amount float64 `json:"amount"`
	}{}

	err := c.BindJSON(&req)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	accId, err := uuid.Parse(req.AccID)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	// build dto
	dto := entity.WithdrawTransactionDto{
		AccID:  accId,
		Desc:   req.Desc,
		Amount: req.Amount,
	}

	// execute usecase
	err = t.TransactionUC.Withdraw(c, &dto)
	if err != nil {
		delivery.HandleError(c, http.StatusInternalServerError)
		return
	}

	delivery.HandleOK(c, nil)
}
