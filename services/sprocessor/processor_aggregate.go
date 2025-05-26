package sprocessor

import (
	"fmt"

	"github.com/TudorHulban/wrpLavinMQ/domain/events"
	"github.com/TudorHulban/wrpLavinMQ/fixtures"
)

// TODO: use sync.Pool for metrics map?
func Agregate(input [][]byte) ([]byte, error) {
	metrics := map[string]*events.MetricInfo{
		fixtures.MetricLabel1: &events.MetricInfo{},
		fixtures.MetricLabel2: &events.MetricInfo{},
		fixtures.MetricLabel3: &events.MetricInfo{},
		fixtures.MetricLabel4: &events.MetricInfo{},
	}

	for _, message := range input {
		event, errCr := events.NewEventA(message)
		if errCr != nil {
			return nil,
				errCr
		}

		values, exists := metrics[event.MetricLabel]
		if !exists {
			return nil,
				fmt.Errorf(
					"procesor agregate could not find metric type: %s",
					event.MetricLabel,
				)
		}

		(*values).NumberOfMessages++
		(*values).Sum = values.Sum + float64(event.Value)
	}

	return events.NewEventB(metrics).AsJSON()
}
