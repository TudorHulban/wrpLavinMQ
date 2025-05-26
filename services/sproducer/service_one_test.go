package sproducer

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/domain/events"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/stretchr/testify/require"
)

func TestProducerOneService(t *testing.T) {
	config, errConfig := configuration.NewConfigurationTest()
	require.NoError(t, errConfig)

	conn, errConnect := connection.Connect(
		&connection.ConfigAMQP{
			Protocol:    config.GetConfigurationValue(configuration.ConfigAMQPProtocol),
			Username:    config.GetConfigurationValue(configuration.ConfigAMQPNameUser),
			Password:    config.GetConfigurationValue(configuration.ConfigAMQPPassword),
			Host:        config.GetConfigurationValue(configuration.ConfigAMQPHost),
			Port:        config.GetConfigurationValue(configuration.ConfigAMQPPort),
			VirtualHost: config.GetConfigurationValue(configuration.ConfigAMQPVirtualHost),
		},
	)
	require.NoError(t, errConnect)
	defer conn.Close()

	require.NotNil(t, conn)

	service := NewServiceProducer(conn)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())

	evA := events.EventA{
		MetricLabel: "jitter",
		Value:       21,
	}

	json, errSerialize := evA.AsJSON()
	require.NoError(t, errSerialize)
	require.NotZero(t, json)

	require.NoError(t,
		service.PublishMessageJSON(
			json,

			&ParamsPublishMessageJSON{
				Exchange: config.GetConfigurationValue(configuration.ConfiqAMQPNameExchange),
				Queue:    config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),
			},
		),
	)
}
