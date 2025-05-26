package sconsumer

import (
	"testing"
	"time"

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

	serviceProcesor, errServiceProcesor := sprocessor.NewServiceProcessor(
		&sprocessor.PiersNewServiceProcessor{
			Configuration: config,
			Proc:          sprocessor.Aggregate,
		},
	)
	require.NoError(t, errServiceProcesor)

	service := NewServiceConsumer(
		&PiersNewServiceConsumer{
			Connection: conn,
			Processor:  serviceProcesor,
		},
	)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())

	go service.Processor.Listen(service.ChProcessorData)

	service.ConsumeContinuoslyMany(
		&ParamsConsume{
			QueueName: config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),

			PefetchCount:              100,
			BatchMaxAggregateDuration: 5 * time.Second,
		},
	)
}
