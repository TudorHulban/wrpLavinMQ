package slogging

import (
	"errors"
	"testing"
)

func TestNewServiceLog(t *testing.T) {
	service := NewServiceLog()
	service.Logger.Info().Msg("some log message")
	service.Logger.Err(errors.New("some error")).Msg("some log message")
}
