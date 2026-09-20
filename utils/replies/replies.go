package replies

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/rs/zerolog/log"
)

func SucceededAnswer(ctx amqp.Context, broker amqp.MessageBroker,
	message *amqp.RabbitMQMessage) {
	reply(ctx, broker, message)
}

// FailedAnswer takes the request so the failure carries its game and language back
// to the requester.
func FailedAnswer(ctx amqp.Context, broker amqp.MessageBroker,
	request *amqp.RabbitMQMessage, messageType amqp.RabbitMQMessage_Type) {
	reply(ctx, broker, amqp.NewFailedReply(request, messageType))
}

func reply(ctx amqp.Context, broker amqp.MessageBroker, message *amqp.RabbitMQMessage) {
	err := broker.Reply(message, ctx.CorrelationID, ctx.ReplyTo)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, ctx.CorrelationID).
			Str(constants.LogGame, message.GetGame().String()).
			Str(constants.LogReplyTo, ctx.ReplyTo).
			Msgf("Cannot publish via broker, request ignored")
	}
}
