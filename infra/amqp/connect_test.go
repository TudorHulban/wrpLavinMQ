package connection

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var _ConfigLocal = ConfigAMQP{
	Protocol: "amqp",
	Username: "guest",
	Password: "guest",
	Host:     "localhost",
	Port:     5672,
}

func TestConnect(t *testing.T) {
	conn, errConnect := Connect(&_ConfigLocal)
	require.NoError(t, errConnect)

	defer conn.Close()

	require.NotNil(t, conn)
}
