package service

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/hdkef/be-assignment/pkg/domain/entity"
	"github.com/hdkef/be-assignment/services/account/domain/service"
	"github.com/streadway/amqp"
)

type PublisherImpl struct {
	Conn                             *amqp.Connection
	publishCreateUserEventsPublisher publishCreateUserEventsPublisher
}

type publishCreateUserEventsPublisher struct {
	Ch *amqp.Channel
}

func (p *publishCreateUserEventsPublisher) publish(ctx context.Context, dto *entity.UserCreatedEventDto) error {
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
		entity.EXCHANGE_USER,           // exchange
		entity.EVENT_NAME_USER_CREATED, // routing key
		false,                          // mandatory
		false,                          // immediate
		msg,                            // message to publish
	)
	if err != nil {
		log.Printf("Error publishing message: %v", err)
		return err
	}

	return nil
}

// PublishCreateUserEvents implements service.Publisher.
func (p *PublisherImpl) PublishCreateUserEvents(ctx context.Context, dto *entity.UserCreatedEventDto) error {
	return p.publishCreateUserEventsPublisher.publish(ctx, dto)
}

func NewPublisher(conn *amqp.Connection) service.Publisher {

	p := &PublisherImpl{
		Conn: conn,
	}

	ch1, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	p.publishCreateUserEventsPublisher = publishCreateUserEventsPublisher{
		Ch: ch1,
	}

	return p
}
