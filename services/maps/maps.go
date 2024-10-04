package maps

import (
	"crypto/rand"
	"fmt"
	"math/big"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker) *Impl {
	return &Impl{
		broker: broker,
	}
}

func (service *Impl) GetMapRequest(request *amqp.CompetitionMapRequest, correlationID,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidGetMapRequest(request) {
		service.publishFailedGetMapAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get competition map request received")

	selectedMap := request.MapNumber
	if selectedMap == 0 {
		n, err := rand.Int(rand.Reader, big.NewInt(constants.MapCount))
		if err != nil {
			log.Error().Err(err).Msg("Failed to randomize map number, returning failed answer")
			service.publishFailedGetMapAnswer(correlationID, answersRoutingkey, lg)
			return
		}
		selectedMap = n.Int64() + 1
	}

	service.publishSucceededGetMapAnswer(correlationID, answersRoutingkey, selectedMap, lg)
}

func (service *Impl) publishSucceededGetMapAnswer(correlationID, answersRoutingkey string,
	mapNumber int64, lg amqp.Language) {
	source := constants.GetMapSource()
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		CompetitionMapAnswer: &amqp.CompetitionMapAnswer{
			MapNumber:    mapNumber,
			MapNormalURL: craftMapImageURL(constants.MapTypeNormal, mapNumber),
			MapTacticURL: craftMapImageURL(constants.MapTypeTactical, mapNumber),
			Source: &amqp.Source{
				Name: source.Name,
				Icon: source.Icon,
				Url:  source.URL,
			},
		},
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer, answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func (service *Impl) publishFailedGetMapAnswer(correlationID, answersRoutingkey string,
	lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER,
		Status:   amqp.RabbitMQMessage_FAILED,
		Language: lg,
	}

	err := service.broker.Publish(&message, amqp.ExchangeAnswer,
		answersRoutingkey, correlationID)
	if err != nil {
		log.Error().Err(err).
			Str(constants.LogCorrelationID, correlationID).
			Msgf("Cannot publish via broker, request ignored")
	}
}

func craftMapImageURL(mapType constants.MapType, number int64) string {
	return fmt.Sprintf(constants.KTArenaMapTemplateURL, mapType, number)
}

func isValidGetMapRequest(request *amqp.CompetitionMapRequest) bool {
	return request != nil
}
