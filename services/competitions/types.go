package competitions

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/services/maps"
)

const (
	requestQueueName   = "competitions-requests"
	requestsRoutingkey = "requests.competitions"
)

type Service interface {
	Consume()
}

type Impl struct {
	broker     amqp.MessageBroker
	mapService maps.Service
}
