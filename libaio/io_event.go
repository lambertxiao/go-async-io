package libaio

import "unsafe"

type IOEvent struct {
	data unsafe.Pointer
	cb   *IOCB
	res  int64
	res2 int64
}

func (e *IOEvent) done() error {
	return nil
}
