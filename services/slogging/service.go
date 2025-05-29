package slogging

import (
	"github.com/phuslu/log"
)

type ServiceLogging struct {
	Logger *log.Logger
}

func NewServiceLog() *ServiceLogging {
	return &ServiceLogging{
		Logger: &log.Logger{
			TimeFormat: "15:04:05",
			Caller:     -1,
			Writer: &log.ConsoleWriter{
				ColorOutput:    true,
				QuoteString:    true,
				EndWithMessage: true,
			},
		},
	}
}
