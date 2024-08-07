package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type ENUM_HISTORY_TRX_TYPE string
type ENUM_HISTORY_STATUS string

const (
	ENUM_HISTORY_TRX_TYPE_POSITIVE ENUM_HISTORY_TRX_TYPE = "POSITIVE"
	ENUM_HISTORY_TRX_TYPE_NEGATIVE ENUM_HISTORY_TRX_TYPE = "NEGATIVE"
	ENUM_HISTORY_STATUS_SUCCESS    ENUM_HISTORY_STATUS   = "SUCCESS"
	ENUM_HISTORY_STATUS_FAILED     ENUM_HISTORY_STATUS   = "FAILED"
)

type History struct {
	Id        uuid.UUID
	ReffNum   uuid.UUID
	CreatedAt time.Time
	AccID     uuid.UUID
	TrxType   ENUM_HISTORY_TRX_TYPE
	Amount    float64
	Status    ENUM_HISTORY_STATUS
	Desc      string
	AccID2    uuid.UUID
}

type GetHistoryDto struct {
	AccID  uuid.UUID
	UserID uuid.UUID
	Page   int
	limit  int
}

func (g *GetHistoryDto) SetLimit(lim int) {
	if lim > 1000 {
		g.limit = 1000
	} else {
		g.limit = lim
	}
}

func (g *GetHistoryDto) GetLimit() int {
	return g.limit
}

func (g *GetHistoryDto) SetPage(page int) {
	if page < 1 {
		g.Page = 1
	} else {
		g.Page = page
	}
}

func (g *GetHistoryDto) Validate() error {

	if g.AccID == uuid.Nil {
		return errors.New("accId is required")
	}

	if g.UserID == uuid.Nil {
		return errors.New("userId is required")
	}

	return nil
}

type GetHistoryFilter struct {
	AccID  uuid.UUID
	UserID uuid.UUID
}

type GetHistoryOptions struct {
	Page  int
	Limit int
}
