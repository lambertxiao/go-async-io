package libaio

import (
	"os"
	"syscall"
	"unsafe"
)

func syscall_iosetup(ioctx int, numEvents int) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IO_SETUP, uintptr(numEvents), uintptr(unsafe.Pointer(&ioctx)), 0)
	if errno != 0 {
		return os.NewSyscallError("IO_SETUP", errno)
	}

	return nil
}

func syscall_destory(ioctx int) error {
	_, _, err := syscall.Syscall(syscall.SYS_IO_DESTROY, uintptr(ioctx), 0, 0)
	if err != 0 {
		return os.NewSyscallError("IO_DESTROY", err)
	}
	return nil
}

func syscall_submit(ioctx int, cb *IOCB) error {
	p := unsafe.Pointer(&cb)
	var len int = 1
	_, _, err := syscall.Syscall(syscall.SYS_IO_SUBMIT, uintptr(ioctx), uintptr(len), uintptr(p))
	if err != 0 {
		return os.NewSyscallError("IO_SUBMIT", err)
	}
	return err
}
