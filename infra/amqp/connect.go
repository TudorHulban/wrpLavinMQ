package connection

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type ConfigAMQP struct {
	Protocol    string
	Username    string
	Password    string
	Host        string
	Port        string
	VirtualHost string
}

func (a ConfigAMQP) String() string {
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s",
		a.Protocol,
		a.Username,
		a.Password,
		a.Host,
		a.Port,
		a.VirtualHost,
	)
}

func Connect(config *ConfigAMQP) (*amqp.Connection, error) {
	url := config.String()

	log.Print(url)

	connection, errConnection := amqp.Dial(url)
	if errConnection != nil {
		return nil,
			errConnection
	}

	return connection,
		nil
}
