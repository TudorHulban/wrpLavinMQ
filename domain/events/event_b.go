package events

import "encoding/json"

type MetricAverage struct {
	MetricLabel  string `json:"metriclabel"`
	AverageValue int    `json:"averagevalue"`
}

type EventB struct {
	Averages []MetricAverage `json:"averages"`
}

func (e EventB) AsJSON() ([]byte, error) {
	return json.Marshal(e)
}
