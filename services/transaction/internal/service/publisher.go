package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/transaction/domain/service"
	"github.com/streadway/amqp"
)

type PublisherImpl struct {
	Conn                                    *amqp.Connection
	publishCreateTransactionEventsPublisher publishCreateTransactionEventsPublisher
}

type publishCreateTransactionEventsPublisher struct {
	Ch *amqp.Channel
}

func (p *publishCreateTransactionEventsPublisher) publish(ctx context.Context, dto *entity.TransactionCreatedEventDto) error {
	// Convert the dto to JSON
	body, err := json.Marshal(dto)
	if err != nil {
		log.Printf("Error marshalling UserCreatedEventDto: %v", err)
		return err
	}

	// Create a message to publish
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
		Timestamp:   time.Now(),
	}

	// Publish the message
	err = p.Ch.Publish(
		entity.EXCHANGE_TRANSACTION,           // exchange
		entity.EVENT_NAME_TRANSACTION_CREATED, // routing key
		false,                                 // mandatory
		false,                                 // immediate
		msg,                                   // message to publish
	)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return err
	}

	return nil
}

// PublishCreateTransactionEvents implements service.Publisher.
func (p *PublisherImpl) PublishCreateTransactionEvents(ctx context.Context, dto *entity.TransactionCreatedEventDto) error {
	return p.publishCreateTransactionEventsPublisher.publish(ctx, dto)
}

func NewPublisher(conn *amqp.Connection) service.Publisher {

	p := &PublisherImpl{
		Conn: conn,
	}

	ch1, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	p.publishCreateTransactionEventsPublisher = publishCreateTransactionEventsPublisher{
		Ch: ch1,
	}

	return p
}
