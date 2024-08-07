package consumer

import (
	"context"
	"encoding/json"

	entity2 "github.com/hdkef/be-assignment/pkg/domain/entity"
)

func (c *ConsumerDelivery) SignUpEvent() {

	ch, err := c.Conn.Channel()
	if err != nil {
		panic(err)
	}
	defer ch.Close()

	// Declare the exchange
	err = ch.ExchangeDeclare(
		entity2.EXCHANGE_USER, // name
		"topic",               // type
		true,                  // durable
		false,                 // auto-deleted
		false,                 // internal
		false,                 // no-wait
		nil,                   // arguments
	)
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"user_created_transaction_queue", // name
		true,                             // durable
		false,                            // delete when unused
		false,                            // exclusive
		false,                            // no-wait
		nil,                              // arguments
	)
	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		q.Name,                          // queue name
		entity2.EVENT_NAME_USER_CREATED, // binding key
		entity2.EXCHANGE_USER,           // exchange
		false,                           // no-wait
		nil,                             // arguments
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

			dto := entity2.UserCreatedEventDto{}
			err = json.Unmarshal(d.Body, &dto)
			if err != nil {
				d.Reject(false)
				continue
			}

			// execute usecase
			err = c.AccBalanceUC.CreateAccountBalance(context.Background(), &dto)
			if err != nil {
				d.Reject(false)
				continue
			}

			d.Ack(false)
		}
	}()
	<-forever
}
