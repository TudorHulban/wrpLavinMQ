package sconsumer

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel
}

func NewServiceProducer(conn *amqp.Connection) *Service {
	return &Service{
		conn: conn,
	}
}

func (s *Service) Connect() error {
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

	EventAsJSON []byte

	Mandatory bool
	Immediate bool
}

func (s *Service) PublishMessageJSON(params *ParamsPublishMessageJSON) error {
	return s.channelAMQP.Publish(
		params.Exchange,
		params.Queue,
		params.Mandatory,
		params.Immediate,

		amqp.Publishing{
			ContentType: "application/json",
			Type:        params.MessageType,
			Body:        params.EventAsJSON,
		},
	)
}

func (s *Service) Close() error {
	return s.conn.Close()
}
