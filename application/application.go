package application

import (
	amqp "github.com/kaellybot/kaelly-amqp"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/kaellybot/kaelly-competition/services/competitions"
	"github.com/kaellybot/kaelly-competition/services/maps"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func New() (*Impl, error) {
	// misc
	broker, err := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		[]amqp.Binding{competitions.GetBinding()})
	if err != nil {
		return nil, err
	}

	// services
	mapService := maps.New(broker)
	competitionService := competitions.New(broker, mapService)

	return &Impl{
		competitionService: competitionService,
		broker:             broker,
	}, nil
}

func (app *Impl) Run() error {
	return app.competitionService.Consume()
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
