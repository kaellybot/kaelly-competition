package maps

import (
	"crypto/rand"
	"math/big"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/kaellybot/kaelly-competition/models/mappers"
	"github.com/kaellybot/kaelly-competition/utils/replies"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker) *Impl {
	return &Impl{
		broker: broker,
	}
}

func (service *Impl) GetMapRequest(ctx amqp.Context, message *amqp.RabbitMQMessage) {
	request := message.GetCompetitionMapRequest()
	if !isValidGetMapRequest(request) {
		replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Str(constants.LogGame, message.GetGame().String()).
		Msgf("Get competition map request received")

	selectedMap := request.GetMapNumber()
	if selectedMap == 0 {
		n, err := rand.Int(rand.Reader, big.NewInt(constants.MapCount))
		if err != nil {
			log.Error().Err(err).Msg("Failed to randomize map number, returning failed answer")
			replies.FailedAnswer(ctx, service.broker, message, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER)
			return
		}
		selectedMap = n.Int64() + 1
	}

	response := mappers.MapGetMapAnswer(message, selectedMap)
	replies.SucceededAnswer(ctx, service.broker, response)
}

func isValidGetMapRequest(request *amqp.CompetitionMapRequest) bool {
	return request != nil
}
