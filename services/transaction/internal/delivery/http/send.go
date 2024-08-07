package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/delivery"
	"github.com/hdkef/be-assignment/services/transaction/domain/entity"
)

func (t *HttpHandler) Send(c *gin.Context) {

	req := struct {
		AccID   string  `json:"accId"`
		Desc    string  `json:"desc"`
		Amount  float64 `json:"amount"`
		ToAccID string  `json:"toAccId"`
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

	toAccId, err := uuid.Parse(req.ToAccID)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	// build dto
	dto := entity.SendTransactionDto{
		AccID:   accId,
		Desc:    req.Desc,
		Amount:  req.Amount,
		ToAccID: toAccId,
	}

	// execute usecase
	err = t.TransactionUC.Send(c, &dto)
	if err != nil {
		delivery.HandleError(c, http.StatusInternalServerError)
		return
	}

	delivery.HandleOK(c, nil)
}
