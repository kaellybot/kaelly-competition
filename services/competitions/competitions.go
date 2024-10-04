package competitions

import (
	"context"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/kaellybot/kaelly-competition/services/maps"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker, mapService maps.Service) *Impl {
	return &Impl{
		broker:     broker,
		mapService: mapService,
	}
}

func GetBinding() amqp.Binding {
	return amqp.Binding{
		Exchange:   amqp.ExchangeRequest,
		RoutingKey: requestsRoutingkey,
		Queue:      requestQueueName,
	}
}

func (service *Impl) Consume() error {
	log.Info().Msgf("Consuming competition requests...")
	return service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(_ context.Context,
	message *amqp.RabbitMQMessage, correlationID string) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.Type {
	case amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST:
		service.mapService.GetMapRequest(message.CompetitionMapRequest, correlationID, answersRoutingkey, message.Language)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Type not recognized, request ignored")
	}
}
