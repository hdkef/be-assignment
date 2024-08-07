package entity

const EXCHANGE_ACCOUNT = "account"
const EVENT_NAME_ACCOUNT_CREATED = "account_created"

type AccountCreatedEventDto struct {
	UserID          string
	AccountID       string
	AccountCurrency string
}
