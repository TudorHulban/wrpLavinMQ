package sconsumer

import (
	"log"
	"time"

	"github.com/TudorHulban/wrpLavinMQ/services/sprocessor"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ServiceConsumer struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel

	processor *sprocessor.ServiceProcessor
	chData    chan [][]byte
}

type PiersNewServiceConsumer struct {
	Connection *amqp.Connection
	Processor  *sprocessor.ServiceProcessor
}

// TODO: piers validation.
func NewServiceConsumer(piers *PiersNewServiceConsumer) *ServiceConsumer {
	return &ServiceConsumer{
		conn:      piers.Connection,
		processor: piers.Processor,
		chData:    make(chan [][]byte),
	}
}

func (s *ServiceConsumer) Connect() error {
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

func (s *ServiceConsumer) ConsumeContinuoslyOne(params *ParamsConsume) error {
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

	blocker := make(chan struct{})

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

	<-blocker

	return nil
}

func (s *ServiceConsumer) ConsumeContinuoslyMany(params *ParamsConsume) error {
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

	blocker := make(chan struct{})

	go func() {
		var batch [][]byte
		timer := time.NewTimer(5 * time.Second)
		defer timer.Stop()

		for {
			select {
			case delivered, opened := <-delivery:
				if !opened {
					// delivery channel closed, send whatever we have
					if len(batch) > 0 {
						s.chData <- batch
					}

					close(s.chData)

					return
				}

				log.Printf("received a message: %s", delivered.Body)
				batch = append(batch, delivered.Body)
				delivered.Ack(false)
				log.Print("message acknowledged")

				// Check if we've reached 100 messages
				if len(batch) >= 100 {
					s.chData <- batch
					batch = nil                  // Reset batch
					timer.Reset(5 * time.Second) // Reset timer for next batch
				}

			case <-timer.C:
				// 5 seconds elapsed, send whatever we have
				if len(batch) > 0 {
					s.chData <- batch
					batch = nil
				}
				timer.Reset(5 * time.Second) // Reset timer for next batch
			}
		}
	}()

	<-blocker

	return nil
}
