package entity

import (
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
