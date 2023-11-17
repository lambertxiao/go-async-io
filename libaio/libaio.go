package libaio

import (
	goaio "go-aio"
	"os"
	"syscall"
	"unsafe"
)

type AIOCtx struct {
	fd    *os.File
	ioctx IOContext
	woff  int
}

type Options struct {
	IODepth int
	Flag    int
	Perm    os.FileMode
}

type IOContext uint

func OpenAIOCtx(fpath string, opts Options) (goaio.IOCtx, error) {
	fd, err := os.OpenFile(fpath, syscall.O_DIRECT|opts.Flag, opts.Perm)
	if err != nil {
		return nil, err
	}

	var ioctx IOContext
	_, _, errno := syscall.Syscall(syscall.SYS_IO_SETUP, uintptr(opts.IODepth), uintptr(unsafe.Pointer(&ioctx)), 0)
	if errno != 0 {
		return nil, os.NewSyscallError("IO_SETUP", err)
	}

	return &AIOCtx{fd: fd, ioctx: ioctx}, nil
}

func (c *AIOCtx) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (c *AIOCtx) WriteAt(p []byte, off int64) (n int, err error) {
	cb := newIOCB(c.fd)
	cb.prepareWrite(p, off)
	err = c.submit(cb)
	if err != nil {
		return 0, err
	}
	return 0, nil
}

func (c *AIOCtx) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (c *AIOCtx) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

func (c *AIOCtx) Close() error {
	_, _, err := syscall.Syscall(syscall.SYS_IO_DESTROY, uintptr(c.ioctx), 0, 0)
	if err != 0 {
		return os.NewSyscallError("IO_DESTROY", err)
	}

	return nil
}
