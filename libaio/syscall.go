package libaio

import (
	"os"
	"syscall"
	"unsafe"
)

func syscall_iosetup(ioctx *uint, numEvents int) error {
	_, _, errno := syscall.Syscall(syscall.SYS_IO_SETUP, uintptr(numEvents), uintptr(unsafe.Pointer(ioctx)), 0)
	if errno != 0 {
		return os.NewSyscallError("IO_SETUP", errno)
	}

	return nil
}

func syscall_destory(ioctx uint) error {
	_, _, err := syscall.Syscall(syscall.SYS_IO_DESTROY, uintptr(ioctx), 0, 0)
	if err != 0 {
		return os.NewSyscallError("IO_DESTROY", err)
	}
	return nil
}

// https://www.man7.org/linux/man-pages/man2/io_submit.2.html
func syscall_submit(ioctx uint, cb *IOCB) error {
	p := unsafe.Pointer(&cb)
	var len int = 1
	_, _, err := syscall.Syscall(syscall.SYS_IO_SUBMIT, uintptr(ioctx), uintptr(len), uintptr(p))
	if err != 0 {
		return os.NewSyscallError("IO_SUBMIT", err)
	}
	return nil
}

// https://www.man7.org/linux/man-pages/man2/io_getevents.2.html
func syscall_getevents(ioctx uint, min_nr, nr int, events []IOEvent, timeout timespec) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}

	p := unsafe.Pointer(&events[0])
	n, _, err := syscall.Syscall6(
		syscall.SYS_IO_GETEVENTS,
		uintptr(ioctx),
		uintptr(min_nr),                   // 最小事件数
		uintptr(nr),                       // 最大事件数
		uintptr(p),                        // 存放接收到的事件
		uintptr(unsafe.Pointer(&timeout)), // 设置超时
		uintptr(0),
	)

	if err != 0 {
		return 0, os.NewSyscallError("IO_GETEVENTS", err)
	}

	return int(n), nil
}
