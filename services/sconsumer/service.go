package sconsumer

import (
	"fmt"
	"time"

	goerrors "github.com/TudorHulban/go-errors"
	"github.com/TudorHulban/wrpLavinMQ/helpers"
	"github.com/TudorHulban/wrpLavinMQ/services/slogging"
	"github.com/TudorHulban/wrpLavinMQ/services/sprocessor"
	"github.com/asaskevich/govalidator"
	amqp "github.com/rabbitmq/amqp091-go"
)

type ServiceConsumer struct {
	conn        *amqp.Connection
	channelAMQP *amqp.Channel

	Processor *sprocessor.ServiceProcessor
	loger     *slogging.ServiceLogging

	ChProcessorData chan [][]byte
}

type PiersNewServiceConsumer struct {
	Connection *amqp.Connection
	Processor  *sprocessor.ServiceProcessor
	Loger      *slogging.ServiceLogging
}

func NewServiceConsumer(piers *PiersNewServiceConsumer) (*ServiceConsumer, error) {
	if errValidatePiers := helpers.ValidatePiers(piers); errValidatePiers != nil {
		return nil,
			goerrors.ErrServiceValidation{
				ServiceName: _ServiceName,
				Caller:      "NewServiceConsumer",
				Issue:       errValidatePiers,
			}
	}

	return &ServiceConsumer{
			conn: piers.Connection,

			Processor: piers.Processor,
			loger:     piers.Loger,

			ChProcessorData: make(chan [][]byte),
		},
		nil
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

	QueueName   string `valid:"required"`
	ConsumerTag string

	BatchMaxAggregateDuration time.Duration `valid:"required"`

	PefetchCount int `valid:"required"`
	PrefetchSize int
	Global       bool

	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
}

func (s *ServiceConsumer) ConsumeContinuoslyMany(params *ParamsConsume) error {
	if _, errValidation := govalidator.ValidateStruct(params); errValidation != nil {
		return goerrors.ErrServiceValidation{
			ServiceName: _ServiceName,
			Caller:      "ConsumeContinuoslyMany",
			Issue:       errValidation,
		}
	}

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

	start := time.Now()

	trackTimeAt := 10000

	var ix int

	trackTime := func(message string, howMany int) {
		s.loger.Logger.Info().Msgf(
			"%s: processed %d messages in %.4fs.\n",
			message,
			howMany,
			time.Since(start).Seconds(),
		)

		start = time.Now()
	}

	go func() {
		var batch [][]byte
		timer := time.NewTimer(params.BatchMaxAggregateDuration)
		defer timer.Stop()

		for {
			select {
			case delivered, opened := <-delivery:
				if !opened {
					if len(batch) > 0 {
						s.ChProcessorData <- batch

						trackTime(
							"delivery channel closed, send whatever we have",
							len(batch),
						)
					}

					close(s.ChProcessorData)

					return
				}

				batch = append(batch, delivered.Body)
				delivered.Ack(false)

				ix++

				if len(batch) >= params.PefetchCount {
					s.ChProcessorData <- batch
					batch = nil

					timer.Reset(params.BatchMaxAggregateDuration)
				}

				if ix%trackTimeAt == 0 {
					go trackTime(
						fmt.Sprintf(
							"reporting at %d messages",
							trackTimeAt,
						),
						ix,
					)

					ix = 0
				}

			case <-timer.C:
				if len(batch) > 0 {
					s.ChProcessorData <- batch

					trackTime(
						"time elapsed, send whatever we have",
						len(batch),
					)

					batch = nil
				}

				timer.Reset(params.BatchMaxAggregateDuration)
			}
		}
	}()

	<-blocker

	return nil
}
