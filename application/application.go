package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/kaellybot/kaelly-competition/services/competitions"
	"github.com/kaellybot/kaelly-competition/services/maps"
	"github.com/kaellybot/kaelly-competition/utils/insights"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(competitions.GetBinding()))
	probes := insights.NewProbes(broker.IsConnected)
	prom := insights.NewPrometheusMetrics()

	// services
	mapService := maps.New(broker)
	competitionService := competitions.New(broker, mapService)

	return &Impl{
		competitionService: competitionService,
		broker:             broker,
		probes:             probes,
		prom:               prom,
	}, nil
}

func (app *Impl) Run() error {
	app.probes.ListenAndServe()
	app.prom.ListenAndServe()

	if err := app.broker.Run(); err != nil {
		return err
	}

	app.competitionService.Consume()
	return nil
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	app.prom.Shutdown()
	app.probes.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
