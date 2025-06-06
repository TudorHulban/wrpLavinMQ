package smetrics

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Metrics map[IdentifierEmitter]float64

func NewMetrics() *Metrics {
	metrics := make(map[IdentifierEmitter]float64)

	return (*Metrics)(&metrics)
}

func (m Metrics) String() string {
	currentTime := time.Now().Format("15:04:05.000")

	var builder strings.Builder

	fmt.Fprintf(
		&builder,
		"Metrics (%s)\n",
		currentTime,
	)

	if len(m) == 0 {
		return builder.String()
	}

	keys := make([]IdentifierEmitter, 0, len(m))

	for k := range m {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	for _, k := range keys {
		fmt.Fprintf(
			&builder,
			"Emitter %d: %.2f\n",
			k,
			m[k],
		)
	}

	return builder.String()
}
