package mappers

import (
	"fmt"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
)

func MapGetMapAnswer(message *amqp.RabbitMQMessage, mapNumber int64) *amqp.RabbitMQMessage {
	source := constants.GetMapSource()
	answer := amqp.NewReply(message, amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER,
		amqp.RabbitMQMessage_SUCCESS)
	answer.CompetitionMapAnswer = &amqp.CompetitionMapAnswer{
		MapNumber:      mapNumber,
		MapNormalURL:   craftMapImageURL(constants.MapTypeNormal, mapNumber),
		MapTacticalURL: craftMapImageURL(constants.MapTypeTactical, mapNumber),
		Source: &amqp.Source{
			Name: source.Name,
			Icon: source.Icon,
			Url:  source.URL,
		},
	}
	return answer
}

func craftMapImageURL(mapType constants.MapType, number int64) string {
	return fmt.Sprintf(constants.KTArenaMapTemplateURL, mapType, number)
}
