package sconsumer

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/domain/events"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/stretchr/testify/require"
)

const (
	_NameExchange = "exchange12345"
	_NameQueue    = "queue12345"
)

func TestService(t *testing.T) {
	config, errConfig := configuration.NewConfigurationTest()
	require.NoError(t, errConfig)

	conn, errConnect := connection.Connect(
		&connection.ConfigAMQP{
			Protocol: config.GetValue(configuration.ConfigAMQPProtocol),
			Username: config.GetValue(configuration.ConfigAMQPUserName),
			Password: config.GetValue(configuration.ConfigAMQPPassword),
			Host:     config.GetValue(configuration.ConfigAMQPHost),
			Port:     config.GetValue(configuration.ConfigAMQPPort),
		},
	)
	require.NoError(t, errConnect)
	defer conn.Close()

	require.NotNil(t, conn)

	service := NewService(conn)
	require.NotNil(t, service)

	require.NoError(t, service.Connect())
	require.NoError(t,
		service.DeclareExchange(
			&ParamsDeclareExchange{
				Name:    _NameExchange,
				Kind:    "direct", // TODO: add type constants.
				Durable: true,
			},
		),
	)
	// require.NoError(t,
	// 	service.DeclareQueue(
	// 		&ParamsDeclareQueue{
	// 			Name: _NameQueue,
	// 		},
	// 	),
	// )

	evA := events.EventA{
		MetricLabel: "jitter",
		Value:       21,
	}

	json, errSerialize := evA.AsJSON()
	require.NoError(t, errSerialize)
	require.NotZero(t, json)

	require.NoError(t,
		service.PublishMessageJSON(
			&ParamsPublishMessageJSON{
				Exchange: _NameExchange,
				Queue:    _NameQueue,

				EventAsJSON: json,
			},
		),
	)
}
