package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/services/competitions"
	"github.com/kaellybot/kaelly-competition/utils/insights"
)

type Application interface {
	Run() error
	Shutdown()
}

type Impl struct {
	competitionService competitions.Service
	broker             amqp.MessageBroker
	probes             insights.Probes
	prom               insights.PrometheusMetrics
}
