package fixtures

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

type IdentifierEmitter uint8

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
