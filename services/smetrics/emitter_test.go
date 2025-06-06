package smetrics

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

type Emitter struct {
	timeInterval   time.Duration
	numberEmitters uint8
}

type Emited struct {
	Identifier IdentifierEmitter
	Payload    float64
}

type ParamsEmitter struct {
	TimeInterval   time.Duration
	NumberEmitters IdentifierEmitter
}

func (e *Emitter) Emit(ctx context.Context) <-chan Emited {
	chEmit := make(chan Emited)

	var wg sync.WaitGroup

	wg.Add(int(e.numberEmitters))

	for ix := range e.numberEmitters {
		go func() {
			defer wg.Done()

			timer := time.NewTimer(e.timeInterval)
			defer timer.Stop()

			for {
				select {
				case <-ctx.Done():
					return

				case <-timer.C:
					chEmit <- Emited{
						Identifier: IdentifierEmitter(ix),
						Payload:    float64(ix) + rand.Float64(),
					}

					timer.Reset(e.timeInterval)
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(chEmit)
	}()

	return chEmit
}

func TestEmitter(t *testing.T) {
	e := Emitter{
		timeInterval:   200 * time.Millisecond,
		numberEmitters: 5,
	}

	values := NewValues(100)

	metrics := NewMetrics()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	chRead := e.Emit(
		ctx,
	)

	startTime := time.Now()
	var messageCount uint

	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	for event := range chRead {
		values.AddValue(
			event.Identifier,
			event.Payload,
		)

		(*metrics)[event.Identifier] = values.GetMetric(event.Identifier)

		messageCount++

		select {
		case <-ticker.C:
			go fmt.Printf(
				"Number values: %d. %s",

				numberValues,
				metrics,
			)

		default:
		}
	}

	elapsedTime := time.Since(startTime).Seconds()
	messagesPerSecond := float64(messageCount) / elapsedTime

	fmt.Printf(
		"Received %d messages in %.2f seconds (%.2f messages/second)\n",
		messageCount,
		elapsedTime,
		messagesPerSecond,
	)
}
