package sprocessor

type Processor func([][]byte) []byte

type ServiceProcessor struct {
	proc        Processor
	ChannelData chan ([][]byte)
}

func NewServiceProcessor(proc Processor) *ServiceProcessor {
	return &ServiceProcessor{
		proc:        proc,
		ChannelData: make(chan [][]byte),
	}
}

func (s *ServiceProcessor) Listen(onChannel chan ([][]byte)) {
	for messages := range onChannel {
		s.proc(messages)
	}
}

// for testing only.
func (s *ServiceProcessor) ListenWithTiming(onChannel chan ([][]byte)) {
	for messages := range onChannel {
		s.proc(messages)
	}
}
