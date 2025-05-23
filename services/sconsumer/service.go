package sconsumer

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel
}

func NewService(conn *amqp.Connection) *Service {
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
