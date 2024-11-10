package mappers

import (
	"fmt"

	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
)

func MapGetMapAnswer(mapNumber int64, lg amqp.Language) *amqp.RabbitMQMessage {
	source := constants.GetMapSource()
	return &amqp.RabbitMQMessage{
		Type:     amqp.RabbitMQMessage_COMPETITION_MAP_ANSWER,
		Status:   amqp.RabbitMQMessage_SUCCESS,
		Language: lg,
		CompetitionMapAnswer: &amqp.CompetitionMapAnswer{
			MapNumber:      mapNumber,
			MapNormalURL:   craftMapImageURL(constants.MapTypeNormal, mapNumber),
			MapTacticalURL: craftMapImageURL(constants.MapTypeTactical, mapNumber),
			Source: &amqp.Source{
				Name: source.Name,
				Icon: source.Icon,
				Url:  source.URL,
			},
		},
	}
}

func craftMapImageURL(mapType constants.MapType, number int64) string {
	return fmt.Sprintf(constants.KTArenaMapTemplateURL, mapType, number)
}
