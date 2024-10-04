package maps

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

type Service interface {
	GetMapRequest(request *amqp.CompetitionMapRequest, correlationID,
		answersRoutingkey string, lg amqp.Language)
}

type Impl struct {
	broker amqp.MessageBroker
}
