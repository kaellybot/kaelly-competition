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

func (service *Impl) GetMapRequest(ctx amqp.Context, request *amqp.CompetitionMapRequest,
	lg amqp.Language) {
	if !isValidGetMapRequest(request) {
		replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, ctx.CorrelationID).
		Msgf("Get competition map request received")

	selectedMap := request.MapNumber
	if selectedMap == 0 {
		n, err := rand.Int(rand.Reader, big.NewInt(constants.MapCount))
		if err != nil {
			log.Error().Err(err).Msg("Failed to randomize map number, returning failed answer")
			replies.FailedAnswer(ctx, service.broker, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER, lg)
			return
		}
		selectedMap = n.Int64() + 1
	}

	response := mappers.MapGetMapAnswer(selectedMap, lg)
	replies.SucceededAnswer(ctx, service.broker, response)
}

func isValidGetMapRequest(request *amqp.CompetitionMapRequest) bool {
	return request != nil
}
