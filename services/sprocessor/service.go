package sprocessor

import (
	"log"

	goerrors "github.com/TudorHulban/go-errors"
	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/helpers"
	"github.com/TudorHulban/wrpLavinMQ/services/sproducer"
)

type Processor func([][]byte) ([][]byte, error)

type ServiceProcessor struct {
	configuration configuration.IConfiguration

	proc        Processor
	ChannelData chan ([][]byte)

	Producer *sproducer.ServiceProducer // TODO: move to io.writer?
}

type PiersNewServiceProcessor struct {
	Configuration configuration.IConfiguration
	Proc          Processor
	Producer      *sproducer.ServiceProducer
}

func NewServiceProcessor(piers *PiersNewServiceProcessor) (*ServiceProcessor, error) {
	if errValidatePiers := helpers.ValidatePiers(piers); errValidatePiers != nil {
		return nil,
			goerrors.ErrServiceValidation{
				ServiceName: _ServiceName,
				Caller:      "NewServiceProcessor",
				Issue:       errValidatePiers,
			}
	}

	return &ServiceProcessor{
			configuration: piers.Configuration,
			proc:          piers.Proc,
			ChannelData:   make(chan [][]byte),
			Producer:      piers.Producer,
		},
		nil
}

func (s *ServiceProcessor) Listen(onChannel chan ([][]byte)) {
	for messages := range onChannel {
		processed, errProcess := s.proc(messages)
		if errProcess != nil {
			go log.Println(errProcess)

			continue
		}

		for _, message := range processed {
			s.Producer.PublishMessageJSON(
				message,

				&sproducer.ParamsPublishMessageJSON{
					Exchange: s.configuration.GetConfigurationValue(configuration.ConfiqAMQPNameExchange),
					Queue:    s.configuration.GetConfigurationValue(configuration.ConfiqAMQPNameQueueAggregates),
				},
			)
		}
	}
}
