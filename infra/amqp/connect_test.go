package connection

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/domain/configurationtest"
	"github.com/stretchr/testify/require"
)

var _ConfigLocal = ConfigAMQP{
	Protocol: "amqp",
	Username: "guest",
	Password: "guest",
	Host:     "localhost",
	Port:     "5672",
}

func TestConnect(t *testing.T) {
	config, errConfig := configurationtest.NewConfigurationTest()
	require.NoError(t, errConfig)

	conn, errConnect := Connect(&_ConfigLocal)
	require.NoError(t, errConnect)

	defer conn.Close()

	require.NotNil(t, conn)
}
