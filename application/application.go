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
	broker := amqp.New(constants.RabbitMQClientID, viper.GetString(constants.RabbitMQAddress),
		amqp.WithBindings(competitions.GetBinding()))

	// services
	mapService := maps.New(broker)
	competitionService := competitions.New(broker, mapService)

	return &Impl{
		competitionService: competitionService,
		broker:             broker,
	}, nil
}

func (app *Impl) Run() error {
	if err := app.broker.Run(); err != nil {
		return err
	}

	app.competitionService.Consume()
	return nil
}

func (app *Impl) Shutdown() {
	app.broker.Shutdown()
	log.Info().Msgf("Application is no longer running")
}
