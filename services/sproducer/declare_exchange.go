package sconsumer

type ParamsDeclareExchange struct {
	Table map[string]any

	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
}

func (s *Service) DeclareExchange(params *ParamsDeclareExchange) error {
	return s.channelAMQP.ExchangeDeclare(
		params.Name,
		params.Kind,
		params.Durable,
		params.AutoDelete,
		params.Internal,
		params.NoWait,
		params.Table,
	)
}
