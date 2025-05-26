package sconsumer

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/stretchr/testify/require"
)

func TestConsumerOneService(t *testing.T) {
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

	service := NewServiceConsumer(
		&PiersNewServiceConsumer{
			Connection: conn,
		},
	)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())

	service.ConsumeContinuoslyOne(
		&ParamsConsume{
			QueueName: config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),
		},
	)
}
