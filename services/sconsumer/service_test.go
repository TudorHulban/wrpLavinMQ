package sconsumer

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/stretchr/testify/require"
)

func TestConsumerService(t *testing.T) {
	config, errConfig := configuration.NewConfigurationTest()
	require.NoError(t, errConfig)

	conn, errConnect := connection.Connect(
		&connection.ConfigAMQP{
			Protocol:    config.GetValue(configuration.ConfigAMQPProtocol),
			Username:    config.GetValue(configuration.ConfigAMQPNameUser),
			Password:    config.GetValue(configuration.ConfigAMQPPassword),
			Host:        config.GetValue(configuration.ConfigAMQPHost),
			Port:        config.GetValue(configuration.ConfigAMQPPort),
			VirtualHost: config.GetValue(configuration.ConfigAMQPVirtualHost),
		},
	)
	require.NoError(t, errConnect)
	defer conn.Close()

	require.NotNil(t, conn)

	service := NewServiceConsumer(conn)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())

	service.ConsumeContinuosly(
		&ParamsConsume{
			QueueName: config.GetValue(configuration.ConfiqAMQPNameQueue),
		},
	)
}
