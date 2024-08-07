package entity

const EXCHANGE_USER = "user"
const EVENT_NAME_USER_CREATED = "user_created"

type UserCreatedEventDto struct {
	UserID          string
	AccountID       string
	AccountCurrency string
}
