package entity

import (
	"time"

	"github.com/google/uuid"
)

type ENUM_SCHEDULE_TRX_STATUS string
type ENUM_SCHEDULE_TYPE string
type ENUM_SCHEDULE_SCHEDULE string

const (
	ENUM_SCHEDULE_TRX_STATUS_ACTIVE   ENUM_SCHEDULE_TRX_STATUS = "ACTIVE"
	ENUM_SCHEDULE_TRX_STATUS_INACTIVE ENUM_SCHEDULE_TRX_STATUS = "INACTIVE"
	ENUM_SCHEDULE_TYPE_SEND           ENUM_SCHEDULE_TYPE       = "SEND"
	ENUM_SCHEDULE_SCHEDULE_DAILY      ENUM_SCHEDULE_SCHEDULE   = "DAILY"
)

type ScheduleTrx struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	Status      ENUM_SCHEDULE_TRX_STATUS
	AccID       uuid.UUID
	Type        ENUM_SCHEDULE_TYPE
	Schedule    ENUM_SCHEDULE_SCHEDULE
	ToAccID     uuid.UUID
	Amount      float64
	HasChecked  bool
	LastChecked *time.Time
}
