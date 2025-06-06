package fixtures

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type Value struct {
	valueNext     *Value
	valuePrevious *Value

	Payload float64
}

type ListExtremities struct {
	Head *Value
	Tail *Value
}

type Values struct {
	values        map[IdentifierEmitter]*ListExtremities
	maxListLength uint16
}

func NewValues(maxListLength uint16) *Values {
	return &Values{
		maxListLength: maxListLength,
		values:        map[IdentifierEmitter]*ListExtremities{},
	}
}

func (v *Values) AddValue(onList IdentifierEmitter, payload float64) {
	if currentExtremities, exists := v.values[onList]; exists {
		newHead := Value{
			valueNext: currentExtremities.Head,
			Payload:   payload,
		}

		if currentExtremities.Head.valueNext == currentExtremities.Tail {
			currentExtremities.Head.valueNext = nil
			currentExtremities.Head.valuePrevious = &newHead

			v.values[onList] = &ListExtremities{
				Head: &newHead,
				Tail: currentExtremities.Head,
			}
		} else {
			v.values[onList] = &ListExtremities{
				Head: &newHead,
				Tail: currentExtremities.Tail,
			}
		}

		return
	}

	v.values[onList] = &ListExtremities{
		Head: &Value{
			Payload: payload,
		},
	}
}

func (v *Values) GetMetric(forList IdentifierEmitter) float64 {
	ix := v.maxListLength

	var sum float64

	head := v.values[forList].Head

	for head.valueNext != nil && ix > 0 {
		sum = sum + head.Payload
		head = head.valueNext

		ix--
	}

	if ix == v.maxListLength {
		return 0
	}

	return sum / float64(v.maxListLength-ix)
}

func TestEmitter(t *testing.T) {
	e := Emitter{
		timeInterval:   10 * time.Millisecond,
		numberEmitters: 30,
	}

	values := NewValues(100)

	metrics := make(map[IdentifierEmitter]float64)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chRead := e.Emit(
		ctx,
	)

	startTime := time.Now()
	var messageCount uint

	for event := range chRead {
		values.AddValue(
			event.Identifier,
			event.Payload,
		)

		metrics[event.Identifier] = values.GetMetric(event.Identifier)

		messageCount++
	}

	elapsedTime := time.Since(startTime).Seconds()
	messagesPerSecond := float64(messageCount) / elapsedTime

	fmt.Printf(
		"Received %d messages in %.2f seconds (%.2f messages/second)\n",
		messageCount,
		elapsedTime,
		messagesPerSecond,
	)

	fmt.Println(
		metrics,
	)
}
