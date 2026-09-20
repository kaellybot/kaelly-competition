package maps

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

type Service interface {
	GetMapRequest(ctx amqp.Context, message *amqp.RabbitMQMessage)
}

type Impl struct {
	broker amqp.MessageBroker
}
