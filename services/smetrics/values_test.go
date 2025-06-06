package smetrics

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValues(t *testing.T) {
	maxNumberItems := 10

	values := NewValues(uint16(maxNumberItems))
	require.NotNil(t, values)
	require.Empty(t, values.values)

	identifier1 := 1
	payload1Identifier1 := .7

	values.AddValue(
		IdentifierEmitter(identifier1),
		payload1Identifier1,
	)

	linkedListIdentifier1 := values.values[IdentifierEmitter(identifier1)]

	require.NotNil(t,
		linkedListIdentifier1,
	)
	require.EqualValues(t,
		1,
		linkedListIdentifier1.NumberValues,
	)

	payload2Identifier1 := .8

	values.AddValue(
		IdentifierEmitter(identifier1),
		payload2Identifier1,
	)
	require.EqualValues(t,
		2,
		linkedListIdentifier1.NumberValues,
	)
}
