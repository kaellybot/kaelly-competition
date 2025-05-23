package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/kaellybot/kaelly-competition/application"
	"github.com/kaellybot/kaelly-competition/models/constants"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	_ "golang.org/x/crypto/x509roots/fallback"
)

func init() {
	initConfig()
	initLog()
}

func initConfig() {
	viper.SetConfigFile(constants.ConfigFileName)

	for configName, defaultValue := range constants.GetDefaultConfigValues() {
		viper.SetDefault(configName, defaultValue)
	}

	err := viper.ReadInConfig()
	if err != nil {
		log.Debug().Str(constants.LogFileName, constants.ConfigFileName).Msgf("Failed to read config file, continue...")
	}

	viper.AutomaticEnv()
}

func initLog() {
	zerolog.SetGlobalLevel(constants.LogLevelFallback)
	zerolog.CallerMarshalFunc = func(_ uintptr, file string, line int) string {
		short := file
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				short = file[i+1:]
				break
			}
		}
		return fmt.Sprintf("%s:%d", short, line)
	}
	log.Logger = log.With().Caller().Logger()

	logLevel, err := zerolog.ParseLevel(viper.GetString(constants.LogLevel))
	if err != nil {
		log.Warn().Err(err).Msgf("Log level not set, continue with %s...", constants.LogLevelFallback)
	} else {
		zerolog.SetGlobalLevel(logLevel)
		log.Debug().Msgf("Logger level set to '%s'", logLevel)
	}
}

func main() {
	app, err := application.New()
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to instantiate application")
	}

	err = app.Run()
	if err != nil {
		log.Fatal().Err(err).Msgf("Shutting down after failing to run application.")
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Info().Msgf("%s v%s is now running. Press CTRL-C to exit.", constants.InternalName, constants.Version)
	<-sc

	log.Info().Msgf("Gracefully shutting down %s...", constants.InternalName)
	app.Shutdown()
}
