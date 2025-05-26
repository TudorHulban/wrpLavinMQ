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
			Protocol: config.GetConfigurationValue(configuration.ConfigAMQPProtocol),
			Username: config.GetConfigurationValue(configuration.ConfigAMQPNameUser),
			Password: config.GetConfigurationValue(configuration.ConfigAMQPPassword),
			Host:     config.GetConfigurationValue(configuration.ConfigAMQPHost),
			Port:     config.GetConfigurationValue(configuration.ConfigAMQPPort),
		},
	)
	require.NoError(t, errConnect)

	defer conn.Close()

	require.NotNil(t, conn)
}
