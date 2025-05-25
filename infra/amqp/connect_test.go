package connection

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/stretchr/testify/require"
)

func TestConnect(t *testing.T) {
	config, errConfig := configuration.NewConfigurationTest()
	require.NoError(t, errConfig)

	conn, errConnect := Connect(
		&ConfigAMQP{
			Protocol: config.GetValue(configuration.ConfigAMQPProtocol),
			Username: config.GetValue(configuration.ConfigAMQPNameUser),
			Password: config.GetValue(configuration.ConfigAMQPPassword),
			Host:     config.GetValue(configuration.ConfigAMQPHost),
			Port:     config.GetValue(configuration.ConfigAMQPPort),
		},
	)
	require.NoError(t, errConnect)

	defer conn.Close()

	require.NotNil(t, conn)
}
