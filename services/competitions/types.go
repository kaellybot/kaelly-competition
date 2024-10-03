package competitions

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/services/maps"
)

const (
	requestQueueName   = "competitions-requests"
	requestsRoutingkey = "requests.competitions"
	answersRoutingkey  = "answers.competitions"
)

type Service interface {
	Consume() error
}

type Impl struct {
	broker     amqp.MessageBroker
	mapService maps.Service
}
