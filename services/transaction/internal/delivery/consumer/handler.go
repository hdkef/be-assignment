package consumer

import (
	"github.com/hdkef/be-assignment/services/transaction/domain/usecase"
	"github.com/streadway/amqp"
)

type ConsumerDelivery struct {
	Conn         *amqp.Connection
	AccBalanceUC usecase.AccountBalanceUC
}
