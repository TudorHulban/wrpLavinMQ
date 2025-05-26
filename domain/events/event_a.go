package events

import "encoding/json"

type EventA struct {
	MetricLabel string `json:"metriclabel"`
	Value       int    `json:"value"`
}

func NewEventA(input []byte) (*EventA, error) {
	var result EventA

	if errUnmarshal := json.Unmarshal(input, &result); errUnmarshal != nil {
		return nil,
			errUnmarshal
	}

	return &result,
		nil
}

func (e EventA) AsJSON() ([]byte, error) {
	return json.Marshal(e)
}
