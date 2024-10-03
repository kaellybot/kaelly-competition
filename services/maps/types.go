package maps

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

type Service interface {
	// TODO change proto type
	GetMapRequest(request *amqp.AlignGetBookRequest, correlationID,
		answersRoutingkey string, lg amqp.Language)
}

type Impl struct {
	broker amqp.MessageBroker
}
