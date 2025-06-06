package fixtures

import (
	"math/rand"

	"github.com/TudorHulban/wrpLavinMQ/domain/events"
)

const (
	MetricLabel1 = "Param1"
	MetricLabel2 = "Param2"
	MetricLabel3 = "Param3"
	MetricLabel4 = "Param4"
)

func ForEventA(howMany uint16) [][]byte {
	result := make([][]byte, int(howMany), howMany)

	for ix := range howMany {
		var metricLabel string

		switch rand.Intn(4) + 1 {
		case 1:
			metricLabel = MetricLabel1

		case 2:
			metricLabel = MetricLabel2

		case 3:
			metricLabel = MetricLabel3

		case 4:
			metricLabel = MetricLabel4
		}

		eventSerialized, _ := events.EventA{
			MetricLabel: metricLabel,
			Value:       1,
		}.
			AsJSON()

		result[ix] = eventSerialized
	}

	return result
}
