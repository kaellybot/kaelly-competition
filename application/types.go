package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/services/competitions"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	competitionService competitions.Service
	broker             amqp.MessageBroker
}
