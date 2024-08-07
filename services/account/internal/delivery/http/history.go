package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/delivery"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func (t *HttpHandler) GetHistory(c *gin.Context) {

	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
	userID := sessionContainer.GetUserID()

	accId, err := uuid.Parse(c.Query("accId"))
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	uId, err := uuid.Parse(userID)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(c.Query("page"))
	limit, _ := strconv.Atoi(c.Query("limit"))

	dto := entity.GetHistoryDto{
		AccID:  accId,
		UserID: uId,
		Page:   page,
	}
	dto.SetLimit(limit)

	history, err := t.AccUsecase.GetHistory(c, &dto)
	if err != nil {
		delivery.HandleError(c, http.StatusInternalServerError)
		return
	}

	// map response
	response := []struct {
		Id        string  `json:"id"`
		ReffNum   string  `json:"reffNum"`
		CreatedAt string  `json:"createdAt"`
		AccID     string  `json:"accId"`
		TrxType   string  `json:"trxType"`
		Amount    float64 `json:"amount"`
		Status    string  `json:"status"`
		Desc      string  `json:"desc"`
		AccID2    string  `json:"accId2"`
	}{}

	for _, v := range history {

		accID2 := ""
		if v.AccID2 != uuid.Nil {
			accID2 = v.AccID2.String()
		}

		response = append(response, struct {
			Id        string  `json:"id"`
			ReffNum   string  `json:"reffNum"`
			CreatedAt string  `json:"createdAt"`
			AccID     string  `json:"accId"`
			TrxType   string  `json:"trxType"`
			Amount    float64 `json:"amount"`
			Status    string  `json:"status"`
			Desc      string  `json:"desc"`
			AccID2    string  `json:"accId2"`
		}{
			Id:        v.Id.String(),
			ReffNum:   v.ReffNum.String(),
			CreatedAt: v.CreatedAt.String(),
			AccID:     v.AccID.String(),
			TrxType:   string(v.TrxType),
			Amount:    v.Amount,
			Status:    string(v.Status),
			Desc:      v.Desc,
			AccID2:    accID2,
		})
	}

	delivery.HandleOK(c, response)
}
