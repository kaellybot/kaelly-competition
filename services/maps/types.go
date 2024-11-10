package maps

import (
	amqp "github.com/kaellybot/kaelly-amqp"
)

type Service interface {
	GetMapRequest(ctx amqp.Context, request *amqp.CompetitionMapRequest,
		lg amqp.Language)
}

type Impl struct {
	broker amqp.MessageBroker
}
