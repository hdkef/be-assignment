package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hdkef/be-assignment/pkg/delivery"
	"github.com/hdkef/be-assignment/services/account/domain/entity"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

func (t *HttpHandler) CreateAccount(c *gin.Context) {

	sessionContainer := session.GetSessionFromRequestContext(c.Request.Context())
	userID := sessionContainer.GetUserID()

	uId, err := uuid.Parse(userID)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	dto := entity.CreateAccountDto{}
	dto.SetUserID(uId)

	err = c.ShouldBindJSON(&dto)
	if err != nil {
		delivery.HandleError(c, http.StatusBadRequest)
		return
	}

	err = t.AccUsecase.CreateAccount(c, &dto)
	if err != nil {
		delivery.HandleError(c, http.StatusInternalServerError)
		return
	}

	delivery.HandleOK(c, nil)
}
