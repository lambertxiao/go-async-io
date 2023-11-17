package libaio

import (
	"os"
	"syscall"
	"unsafe"
)

func (c *AIOCtx) submit(cb *IOCB) error {
	p := unsafe.Pointer(&cb)
	var len int = 1
	_, _, err := syscall.Syscall(syscall.SYS_IO_SUBMIT, uintptr(c.ioctx), uintptr(len), uintptr(p))
	if err != 0 {
		return os.NewSyscallError("IO_SUBMIT", err)
	}
	return nil
}
