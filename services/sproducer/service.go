package sproducer

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type ServiceProducer struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel
}

func NewServiceProducer(conn *amqp.Connection) *ServiceProducer {
	return &ServiceProducer{
		conn: conn,
	}
}

func (s *ServiceProducer) Connect() error {
	ch, errChannelOpen := s.conn.Channel()
	if errChannelOpen != nil {
		return errChannelOpen
	}

	s.channelAMQP = ch

	return nil
}

type ParamsPublishMessageJSON struct {
	Exchange    string
	Queue       string
	MessageType string

	Mandatory bool
	Immediate bool
}

func (s *ServiceProducer) PublishMessageJSON(message []byte, params *ParamsPublishMessageJSON) error {
	return s.channelAMQP.Publish(
		params.Exchange,
		params.Queue,
		params.Mandatory,
		params.Immediate,

		amqp.Publishing{
			ContentType: "application/json",
			Type:        params.MessageType,
			Body:        message,
		},
	)
}

func (s *ServiceProducer) Close() error {
	return s.conn.Close()
}
