package consumer

import (
	"context"
	"encoding/json"

	entity2 "github.com/hdkef/be-assignment/pkg/domain/entity"
)

func (c *ConsumerDelivery) TransactionCreatedEvent() {

	ch, err := c.Conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
		entity2.EXCHANGE_TRANSACTION, // name
		"topic",                      // type
		true,                         // durable
		false,                        // auto-deleted
		false,                        // internal
		false,                        // no-wait
		nil,                          // arguments
	)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"transaction_created_account_queue", // name
		true,                                // durable
		false,                               // delete when unused
		false,                               // exclusive
		false,                               // no-wait
		nil,                                 // arguments
	)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		q.Name,                                 // queue name
		entity2.EVENT_NAME_TRANSACTION_CREATED, // binding key
		entity2.EXCHANGE_TRANSACTION,           // exchange
		false,                                  // no-wait
		nil,                                    // arguments
	)
	if err != nil {
		panic(err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {

			dto := entity2.TransactionCreatedEventDto{}
			err = json.Unmarshal(d.Body, &dto)
			if err != nil {
				d.Reject(false)
				continue
			}

			// execute usecase
			var err error
			switch dto.EventType {
			case entity2.ENUM_TRX_CREATED_EVENT_TYPE_SEND:
				err = c.AccountUC.TransactionCreatedSend(context.Background(), dto)
			case entity2.ENUM_TRX_CREATED_EVENT_TYPE_WITHDRAW:
				err = c.AccountUC.TransactionCreatedWithdraw(context.Background(), dto)
			}
			if err != nil {
				d.Reject(false)
				continue
			}

			d.Ack(false)
		}
	}()
	<-forever
}
