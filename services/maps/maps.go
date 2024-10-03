package maps

import (
	"fmt"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/rs/zerolog/log"
)

func New(broker amqp.MessageBroker) *Impl {
	return &Impl{
		broker: broker,
	}
}

func (service *Impl) GetMapRequest(request *amqp.AlignGetBookRequest, correlationID,
	answersRoutingkey string, lg amqp.Language) {
	if !isValidGetMapRequest(request) {
		service.publishFailedGetMapAnswer(correlationID, answersRoutingkey, lg)
		return
	}

	log.Info().Str(constants.LogCorrelationID, correlationID).
		Msgf("Get competition map request received")

	//TODO make it work (if 0 => random, else returns right map)

	service.publishSucceededGetMapAnswer(correlationID, answersRoutingkey, lg)
}

func (service *Impl) publishSucceededGetMapAnswer(correlationID, answersRoutingkey string,
	lg amqp.Language) {
	message := amqp.RabbitMQMessage{
		// TODO change proto type
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		// TODO set proto response
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
		// TODO change proto type
		Type:     amqp.RabbitMQMessage_ALIGN_GET_BOOK_ANSWER,
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

func craftMapImageURL(mapType constants.MapType, number int) string {
	return fmt.Sprintf(constants.KTArenaMapTemplateURL, mapType, number)
}

// TODO change proto type
func isValidGetMapRequest(request *amqp.AlignGetBookRequest) bool {
	return request != nil
}
