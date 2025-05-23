package connection

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfigAMQP struct {
	Protocol string
	Username string
	Password string
	Host     string
	Port     string
}

func (a ConfigAMQP) String() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/",
		a.Protocol,
		a.Username,
		a.Password,
		a.Host,
		a.Port,
	)
}

func Connect(config *ConfigAMQP) (*amqp.Connection, error) {
	connection, errConnection := amqp.Dial(config.String())
	if errConnection != nil {
		return nil,
			errConnection
	}

	return connection,
		nil
}
