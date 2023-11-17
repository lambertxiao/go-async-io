package libaio

import "unsafe"

// https://www.man7.org/linux/man-pages/man2/io_getevents.2.html
type IOEvent struct {
	data unsafe.Pointer
	obj  *IOCB
	res  int64
	res2 int64
}

func (e *IOEvent) done() error {
	return nil
}
