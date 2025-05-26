package events

import "github.com/bytedance/sonic"

type MetricInfo struct {
	Sum              float64
	NumberOfMessages int `json:"numbermessages"`
}

func (m *MetricInfo) Average() float64 {
	if m.NumberOfMessages == 0 {
		return 0
	}

	return m.Sum / float64(m.NumberOfMessages)
}

type MetricAverage struct {
	MetricLabel      string  `json:"metriclabel"`
	AverageValue     float64 `json:"averagevalue"`
	NumberOfMessages int     `json:"numbermessages"`
}

type EventB struct {
	Averages []MetricAverage `json:"averages"`
}

func NewEventB(input map[string]*MetricInfo) *EventB {
	var result EventB

	for metric, values := range input {
		result.Averages = append(
			result.Averages,

			MetricAverage{
				MetricLabel:      metric,
				AverageValue:     values.Average(),
				NumberOfMessages: values.NumberOfMessages,
			},
		)
	}

	return &result
}

func (e EventB) AsJSON() ([]byte, error) {
	return sonic.Marshal(e)
}
