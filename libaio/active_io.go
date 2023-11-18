package libaio

import "sync"

// 用于绑定IO和对应的CB
type ActiveIO struct {
	cb       *IOCB
	state    sync.WaitGroup
	retBytes int64
}

func newActiveIO(cb *IOCB) *ActiveIO {
	s := &ActiveIO{cb: cb}
	s.state.Add(1)
	return s
}

func (s *ActiveIO) Done() {
	s.state.Done()
}

func (s *ActiveIO) Wait() {
	s.state.Wait()
}
