package main

import (
	"os"
	"time"

	goerrors "github.com/TudorHulban/go-errors"
	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/fixtures"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/TudorHulban/wrpLavinMQ/services/sconsumer"
	"github.com/TudorHulban/wrpLavinMQ/services/slogging"
	"github.com/TudorHulban/wrpLavinMQ/services/sprocessor"
	"github.com/TudorHulban/wrpLavinMQ/services/sproducer"
)

func main() {
	serviceLoger := slogging.NewServiceLog()

	config, errConfig := configuration.NewConfigurationTest()
	if errConfig != nil {
		serviceLoger.Logger.
			Err(errConfig).
			Msg("configuration setup")
		os.Exit(
			goerrors.OSExitForConfigurationIssues,
		)
	}

	conn, errConnect := connection.Connect(
		&connection.ConfigAMQP{
			Protocol:    config.GetConfigurationValue(configuration.ConfigAMQPProtocol),
			Username:    config.GetConfigurationValue(configuration.ConfigAMQPNameUser),
			Password:    config.GetConfigurationValue(configuration.ConfigAMQPPassword),
			Host:        config.GetConfigurationValue(configuration.ConfigAMQPHost),
			Port:        config.GetConfigurationValue(configuration.ConfigAMQPPort),
			VirtualHost: config.GetConfigurationValue(configuration.ConfigAMQPVirtualHost),
		},
	)
	if errConnect != nil {
		serviceLoger.Logger.
			Err(errConnect).
			Msg("AMQP setup")
		os.Exit(
			goerrors.OSExitForApplicationIssues,
		)
	}
	defer conn.Close()

	serviceLoger.Logger.Info().Msgf(
		"connected to MQ virtual host: %s/%s",
		config.GetConfigurationValue(configuration.ConfigAMQPHost),
		config.GetConfigurationValue(configuration.ConfigAMQPVirtualHost),
	)

	serviceProducer := sproducer.NewServiceProducer(conn)
	if errConnectProducer := serviceProducer.Connect(); errConnectProducer != nil {
		serviceLoger.Logger.
			Err(errConnectProducer).
			Msg("service producer")
		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}

	howMany := 10000

	messages := fixtures.ForEventA(uint16(howMany))

	startTimeProduce := time.Now()

	for _, msg := range messages {
		if errPublish := serviceProducer.PublishMessageJSON(
			msg,

			&sproducer.ParamsPublishMessageJSON{
				Exchange: config.GetConfigurationValue(configuration.ConfiqAMQPNameExchange),
				Queue:    config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),
			},
		); errPublish != nil {
			serviceLoger.Logger.
				Err(errPublish).
				Msg("service producer")
			os.Exit(
				goerrors.OSExitForServiceIssues,
			)
		}
	}

	serviceLoger.Logger.Info().Msgf(
		"produced %d messages in %s\n",
		howMany,
		time.Since(startTimeProduce),
	)

	serviceProcesor, errServiceProcesor := sprocessor.NewServiceProcessor(
		&sprocessor.PiersNewServiceProcessor{
			Configuration: config,
			Proc:          sprocessor.Aggregate,
			// Proc:     sprocessor.PassThrough,
			Producer: serviceProducer,
		},
	)
	if errServiceProcesor != nil {
		serviceLoger.Logger.
			Err(errServiceProcesor).
			Msg("service processor")
		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}

	serviceConsumer, errCrServiceConsummer := sconsumer.NewServiceConsumer(
		&sconsumer.PiersNewServiceConsumer{
			Connection: conn,
			Processor:  serviceProcesor,
			Loger:      serviceLoger,
		},
	)
	if errCrServiceConsummer != nil {
		serviceLoger.Logger.
			Err(errCrServiceConsummer).
			Msg("service consumer")
		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}
	if errConnectServiceConsumer := serviceConsumer.Connect(); errConnectServiceConsumer != nil {
		serviceLoger.Logger.
			Err(errServiceProcesor).
			Msg("service consumer")
		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}

	go serviceConsumer.Processor.ListenSequential(serviceConsumer.ChProcessorData)

	serviceConsumer.ConsumeContinuoslyMany(
		&sconsumer.ParamsConsume{
			QueueName: config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),

			PefetchCount:              100,
			BatchMaxAggregateDuration: 3 * time.Second,
		},
	)
}
