package sproducer

import (
	"fmt"
	"testing"
	"time"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/fixtures"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/stretchr/testify/require"
)

func TestProducerManyService(t *testing.T) {
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

	howMany := 10000

	messages := fixtures.ForEventA(uint16(howMany))

	startTime := time.Now()

	for _, msg := range messages {
		require.NoError(t,
			service.PublishMessageJSON(
				msg,
				&ParamsPublishMessageJSON{
					Exchange: config.GetConfigurationValue(configuration.ConfiqAMQPNameExchange),
					Queue:    config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),
				},
			),
		)
	}

	fmt.Println(
		time.Since(startTime), // 670.355159ms for 10k.
	)
}
