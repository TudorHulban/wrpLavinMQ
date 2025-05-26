package helpers

import (
	"testing"

	"github.com/TudorHulban/wrpLavinMQ/configuration"
	"github.com/stretchr/testify/require"
)

func TestErrorsValidatePiers(t *testing.T) {
	type Object struct {
		Field configuration.IConfiguration
	}

	object := &Object{}

	require.Error(t, ValidatePiers(object))
}
