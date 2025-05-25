package sconsumer

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel
}

func NewServiceConsumer(conn *amqp.Connection) *Service {
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

type ParamsConsume struct {
	Table amqp.Table

	QueueName   string
	ConsumerTag string

	PefetchCount int
	PrefetchSize int
	Global       bool

	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

func (s *Service) ConsumeContinuosly(params *ParamsConsume) error {
	if errQOS := s.channelAMQP.Qos(
		params.PefetchCount,
		params.PrefetchSize,
		params.Global,
	); errQOS != nil {
		return errQOS
	}

	delivery, errConsume := s.channelAMQP.Consume(
		params.QueueName,
		params.ConsumerTag,
		params.AutoAck,
		params.Exclusive,
		params.NoLocal,
		params.NoWait,
		params.Table,
	)
	if errConsume != nil {
		return errConsume
	}

	forever := make(chan struct{})

	go func() {
		for delivered := range delivery {
			log.Printf(
				"received a message: %s",
				delivered.Body,
			)

			// simulate some work
			time.Sleep(2 * time.Second)

			delivered.Ack(false)
			log.Print("message acknowledged")
		}
	}()

	<-forever

	return nil
}
