package constants

import "github.com/rs/zerolog"

const (
	ConfigFileName = ".env"

	// RabbitMQ address.
	RabbitMQAddress = "RABBITMQ_ADDRESS"

	// Probe port.
	ProbePort = "PROBE_PORT"

	// Metric port.
	MetricPort = "METRIC_PORT"

	// Zerolog values from [trace, debug, info, warn, error, fatal, panic].
	LogLevel = "LOG_LEVEL"

	// Boolean; used to register commands at development guild level or globally.
	Production = "PRODUCTION"

	defaultRabbitMQAddress = "amqp://localhost:5672"
	defaultProbePort       = 9090
	defaultMetricPort      = 2112
	defaultLogLevel        = zerolog.InfoLevel
	defaultProduction      = false
)

func GetDefaultConfigValues() map[string]any {
	return map[string]any{
		RabbitMQAddress: defaultRabbitMQAddress,
		ProbePort:       defaultProbePort,
		MetricPort:      defaultMetricPort,
		LogLevel:        defaultLogLevel.String(),
		Production:      defaultProduction,
	}
}
