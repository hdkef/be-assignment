package entity

import (
	"time"

	"github.com/google/uuid"
)

type ENUM_TRX_LOGS_TRX_TYPE string
type ENUM_TRX_LOGS_STATUS string
type ENUM_TRX_LOGS_EVENT_TYPE string

const (
	ENUM_TRX_LOGS_STATUS_SUCCESS      ENUM_TRX_LOGS_STATUS     = "SUCCESS"
	ENUM_TRX_LOGS_STATUS_FAILED       ENUM_TRX_LOGS_STATUS     = "FAILED"
	ENUM_TRX_LOGS_EVENT_TYPE_SEND     ENUM_TRX_LOGS_EVENT_TYPE = "SEND"
	ENUM_TRX_LOGS_EVENT_TYPE_WITHDRAW ENUM_TRX_LOGS_EVENT_TYPE = "WITHDRAW"
)

type TransactionLogs struct {
	ReffNum   uuid.UUID
	CreatedAt time.Time
	AccID     uuid.UUID
	Amount    float64
	Status    ENUM_TRX_LOGS_STATUS
	Desc      string
	ToAccID   *uuid.UUID
	EventType ENUM_TRX_LOGS_EVENT_TYPE
}
