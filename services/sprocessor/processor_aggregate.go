package sprocessor

import (
	"fmt"
	"sync"

	"github.com/TudorHulban/wrpLavinMQ/domain/events"
	"github.com/TudorHulban/wrpLavinMQ/fixtures"
)

var (
	metricsPool = sync.Pool{
		New: func() any {
			return make(map[string]*events.MetricInfo, 4)
		},
	}
)

func Aggregate(input [][]byte) ([]byte, error) {
	metrics := metricsPool.Get().(map[string]*events.MetricInfo)
	defer metricsPool.Put(metrics)

	metrics[fixtures.MetricLabel1] = &events.MetricInfo{}
	metrics[fixtures.MetricLabel2] = &events.MetricInfo{}
	metrics[fixtures.MetricLabel3] = &events.MetricInfo{}
	metrics[fixtures.MetricLabel4] = &events.MetricInfo{}

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

		values.NumberOfMessages++
		values.Sum = values.Sum + float64(event.Value)
	}

	return events.NewEventB(metrics).AsJSON()
}
