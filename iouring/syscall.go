package iouring

import (
	"os"
	"syscall"
	"unsafe"

	"golang.org/x/sys/unix"
)

// Passed in for io_uring_setup(2). Copied back with updated info on success
type IOUringParams struct {
	sq_entries     uint32
	cq_entries     uint32
	flags          uint32
	sq_thread_cpu  uint32
	sq_thread_idle uint32
	features       uint32
	wq_fd          uint32
	resv           [3]uint32
	sq_off         IOSqringOffsets
	cq_off         IOCqringOffsets
}

type IOSqringOffsets struct {
	head         uint32
	tail         uint32
	ring_mask    uint32
	ring_entries uint32
	flags        uint32
	dropped      uint32
	array        uint32
	resv         [3]uint32
}

type IOCqringOffsets struct {
	head         uint32
	tail         uint32
	ring_mask    uint32
	ring_entries uint32
	overflow     uint32
	cqes         uint32
	flags        uint32
	resv         [3]uint32
}

// https://www.man7.org/linux/man-pages/man2/io_uring_setup.2.html
// setup a context for performing asynchronous I/O
func syscall_io_uring_setup(entries uint, params *IOUringParams) (int, error) {
	res, _, errno := syscall.Syscall(
		unix.SYS_IO_URING_SETUP,
		uintptr(entries),
		uintptr(unsafe.Pointer(params)),
		0,
	)
	if errno != 0 {
		return int(res), os.NewSyscallError("io_uring_setup", errno)
	}

	return int(res), nil
}

// https://www.man7.org/linux/man-pages/man2/io_uring_register.2.html
// register files or user buffers for asynchronous I/O
func syscall_io_uring_register(fd int, opcode uint8, args unsafe.Pointer, nrArgs uint32) error {
	for {
		_, _, errno := syscall.Syscall6(
			unix.SYS_IO_URING_REGISTER,
			uintptr(fd),
			uintptr(opcode),
			uintptr(args),
			uintptr(nrArgs),
			0,
			0,
		)
		if errno != 0 {
			if errno == syscall.EINTR {
				continue
			}
			return os.NewSyscallError("io_uring_register", errno)
		}
		return nil
	}
}

// https://www.man7.org/linux/man-pages/man2/io_uring_enter2.2.html
// initiate and/or complete asynchronous I/O
func syscall_io_uring_enter(fd int, toSubmit uint32, minComplete uint32, flags uint32, sigset *unix.Sigset_t) (int, error) {
	res, _, errno := syscall.Syscall6(
		unix.SYS_IO_URING_ENTER,
		uintptr(fd),
		uintptr(toSubmit),
		uintptr(minComplete),
		uintptr(flags),
		uintptr(unsafe.Pointer(sigset)),
		0,
	)
	if errno != 0 {
		return 0, os.NewSyscallError("io_uring_enter", errno)
	}
	if res < 0 {
		return 0, os.NewSyscallError("io_uring_enter", syscall.Errno(-res))
	}

	return int(res), nil
}

func syscall_mmap(fd int, length uint32, offset uint64) (uintptr, error) {
	ptr, _, errno := syscall.Syscall6(
		syscall.SYS_MMAP,
		0,
		uintptr(length),
		syscall.PROT_READ|syscall.PROT_WRITE, // 映射区域的保护标志。这里表示该区域可读也可写。
		syscall.MAP_SHARED|syscall.MAP_POPULATE, // 映射区域的共享类型。这里表示该区域是共享的，并且使用预读（POPULATE）策略。
		uintptr(fd),
		uintptr(offset),
	)
	if errno != 0 {
		return 0, os.NewSyscallError("mmap", errno)
	}
	return uintptr(ptr), nil
}
