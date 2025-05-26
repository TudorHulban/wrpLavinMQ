package sprocessor

import "log"

func Summary(input [][]byte) []byte {
	log.Printf(
		"number messages: %d",
		len(input),
	)

	return nil
}
