package main

import (
	"fmt"
	"os"
	"time"

	goerrors "github.com/TudorHulban/go-errors"
	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/TudorHulban/wrpLavinMQ/fixtures"
	connection "github.com/TudorHulban/wrpLavinMQ/infra/amqp"
	"github.com/TudorHulban/wrpLavinMQ/services/sconsumer"
	"github.com/TudorHulban/wrpLavinMQ/services/sprocessor"
	"github.com/TudorHulban/wrpLavinMQ/services/sproducer"
)

func main() {
	config, errConfig := configuration.NewConfigurationTest()
	if errConfig != nil {
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
		os.Exit(
			goerrors.OSExitForApplicationIssues,
		)
	}
	defer conn.Close()

	serviceProducer := sproducer.NewServiceProducer(conn)
	if errConnectProducer := serviceProducer.Connect(); errConnectProducer != nil {
		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}
	if errConnectServiceProducer := serviceProducer.Connect(); errConnectServiceProducer != nil {
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
			fmt.Println(errPublish)

			os.Exit(
				goerrors.OSExitForServiceIssues,
			)
		}
	}

	fmt.Printf(
		"produced %d messages in %s\n",
		howMany,
		time.Since(startTimeProduce),
	)

	serviceProcesor, errServiceProcesor := sprocessor.NewServiceProcessor(
		&sprocessor.PiersNewServiceProcessor{
			Configuration: config,
			// Proc:          sprocessor.Aggregate,
			Proc:     sprocessor.PassThrough,
			Producer: serviceProducer,
		},
	)
	if errServiceProcesor != nil {
		fmt.Println(errServiceProcesor)

		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}

	serviceConsumer := sconsumer.NewServiceConsumer(
		&sconsumer.PiersNewServiceConsumer{
			Connection: conn,
			Processor:  serviceProcesor,
		},
	)
	if errConnectServiceConsumer := serviceConsumer.Connect(); errConnectServiceConsumer != nil {
		fmt.Println(errConnectServiceConsumer)

		os.Exit(
			goerrors.OSExitForServiceIssues,
		)
	}

	go serviceConsumer.Processor.Listen(serviceConsumer.ChProcessorData)

	serviceConsumer.ConsumeContinuoslyMany(
		&sconsumer.ParamsConsume{
			QueueName: config.GetConfigurationValue(configuration.ConfiqAMQPNameQueueMessages),

			PefetchCount:              100,
			BatchMaxAggregateDuration: 3 * time.Second,
		},
	)
}
