package sconsumer

import "fmt"

type ParamsDeclareQueue struct {
	Table map[string]any

	Name           string
	BindToExchange string
	Durable        bool
	AutoDelete     bool
	Exclusive      bool
	NoWait         bool
}

func (s *Service) DeclareQueue(params *ParamsDeclareQueue) error {
	_, errDeclare := s.channelAMQP.QueueDeclare(
		params.Name,
		params.Durable,
		params.AutoDelete,
		params.Exclusive,
		params.NoWait,
		params.Table,
	)

	if errDeclare != nil {
		return errDeclare
	}

	errBind := s.channelAMQP.QueueBind(
		params.Name,
		params.Name,
		params.BindToExchange,
		params.NoWait,
		params.Table,
	)

	return fmt.Errorf(
		"queue: %s, bind error: %w",
		params.Name,
		errBind,
	)
}
