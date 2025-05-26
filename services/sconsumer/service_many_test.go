package sconsumer

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/TudorHulban/wrpLavinMQ/services/sprocessor"
	"github.com/stretchr/testify/require"
)

func TestConsumerManyService(t *testing.T) {
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
			Processor:  sprocessor.NewServiceProcessor(sprocessor.Summary),
		},
	)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())

	go service.processor.Listen(service.chData)

	service.ConsumeContinuoslyMany(
		&ParamsConsume{
			QueueName: config.GetConfigurationValue(configuration.ConfiqAMQPNameQueue),

			PefetchCount: 100,
		},
	)
}
