package competitions

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/kaellybot/kaelly-competition/services/maps"
	"github.com/kaellybot/kaelly-competition/utils/replies"
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

func (service *Impl) Consume() {
	log.Info().Msgf("Consuming competition requests...")
	service.broker.Consume(requestQueueName, service.consume)
}

func (service *Impl) consume(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	//exhaustive:ignore Don't need to be exhaustive here since they will be handled by default case
	switch message.GetType() {
	case amqp.RabbitMQMessage_COMPETITION_MAP_REQUEST:
		service.handle(ctx, message, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER,
			service.mapService.GetMapRequest)
	default:
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Type not recognized, request ignored")
	}
}

// handle validates the request game before running the handler, so that a request
// never gets another game's data.
func (service *Impl) handle(ctx amqp.Context, message *amqp.RabbitMQMessage,
	answerType amqp.RabbitMQMessage_Type, handler requestHandler) {
	// The broker guards already refuse these, so reaching this point means the guard
	// was bypassed. Answering is impossible: the reply would carry ANY_GAME too.
	if !message.IsGameSet() {
		log.Error().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Msgf("Request without game received, request ignored")
		return
	}

	if !isGameSupported(message.GetGame()) {
		log.Warn().
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogGame, message.GetGame().String()).
			Msgf("Game not supported, request refused")
		replies.FailedAnswer(ctx, service.broker, message, answerType)
		return
	}

	handler(ctx, message)
}

// isGameSupported reports whether competition can serve this game. KTArena maps are
// Dofus-only, so every other game is refused rather than served Dofus data.
func isGameSupported(game amqp.Game) bool {
	switch game {
	case amqp.Game_DOFUS_GAME:
		return true
	case amqp.Game_ANY_GAME, amqp.Game_DOFUS_TOUCH, amqp.Game_DOFUS_RETRO:
		return false
	default:
		return false
	}
}
