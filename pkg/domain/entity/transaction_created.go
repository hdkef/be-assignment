package entity

import "time"

const EXCHANGE_TRANSACTION = "transaction"
const EVENT_NAME_TRANSACTION_CREATED = "transaction.created"

type ENUM_TRX_CREATED_STATUS string
type ENUM_TRX_CREATED_EVENT_TYPE string

const (
	ENUM_TRX_CREATED_STATUS_SUCCESS      ENUM_TRX_CREATED_STATUS     = "SUCCESS"
	ENUM_TRX_CREATED_STATUS_FAILED       ENUM_TRX_CREATED_STATUS     = "FAILED"
	ENUM_TRX_CREATED_EVENT_TYPE_SEND     ENUM_TRX_CREATED_EVENT_TYPE = "SEND"
	ENUM_TRX_CREATED_EVENT_TYPE_WITHDRAW ENUM_TRX_CREATED_EVENT_TYPE = "WITHDRAW"
)

type TransactionCreatedEventDtoDetail struct {
	CreatedAt time.Time
	AccID     string
	Amount    float64
	Status    ENUM_TRX_CREATED_STATUS
	Desc      string
	ToAccID   *string
}

type TransactionCreatedEventDto struct {
	ReffNum   string
	EventType ENUM_TRX_CREATED_EVENT_TYPE
	Detail    TransactionCreatedEventDtoDetail
}
