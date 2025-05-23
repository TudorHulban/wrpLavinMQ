package events

import "encoding/json"

type EventA struct {
	MetricLabel string `json:"metriclabel"`
	Value       int    `json:"value"`
}

func (e EventA) AsJSON() ([]byte, error) {
	return json.Marshal(e)
}
