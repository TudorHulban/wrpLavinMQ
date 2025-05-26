package fixtures

import (
	"math/rand"

	"github.com/TudorHulban/wrpLavinMQ/domain/events"
)

const (
	_MetricLabel1 = "Param1"
	_MetricLabel2 = "Param2"
	_MetricLabel3 = "Param3"
	_MetricLabel4 = "Param4"
)

func ForEventA(howMany uint16) [][]byte {
	result := make([][]byte, int(howMany), howMany)

	for ix := range howMany {
		var metricLabel string

		switch rand.Intn(4) + 1 {
		case 1:
			metricLabel = _MetricLabel1

		case 2:
			metricLabel = _MetricLabel1

		case 3:
			metricLabel = _MetricLabel1

		case 4:
			metricLabel = _MetricLabel1
		}

		eventSerialized, _ := events.EventA{
			MetricLabel: metricLabel,
			Value:       1,
		}.AsJSON()

		result[ix] = eventSerialized
	}

	return result
}
