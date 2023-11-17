package libaio

import (
	goaio "go-aio"
	"os"
	"syscall"
	"unsafe"

	"github.com/hashicorp/go-multierror"
)

type AIOCtx struct {
	fd      *os.File
	ioctx   IOContext // 由内核libaio填充该iocxt
	woff    int64
	events  []IOEvent
	timeout timespec
}

type Options struct {
	IODepth int
	Flag    int
	Perm    os.FileMode
	// libaio timeout, unit ms 0 means no timeout.
	Timeout int
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

	t := timespec{
		sec:  opts.Timeout / 1000,
		nsec: (opts.Timeout % 1000) * 1000 * 1000,
	}
	events := make([]IOEvent, opts.IODepth)

	ctx := &AIOCtx{
		fd:      fd,
		ioctx:   ioctx,
		events:  events,
		timeout: t,
	}

	go ctx.loop()
	return ctx, nil
}

func (c *AIOCtx) loop() {
	for {
		c.waitEvents()
	}
}

func (c *AIOCtx) waitEvents() error {
	n, err := c.syscall_GetEvents(1, 1, c.events, c.timeout)
	if err != nil {
		return err
	}

	var errs error
	for i := 0; i < n; i++ {
		err := c.events[i].done()
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

// 获取至少从min_nr到nr范围内的事件个数
func (c *AIOCtx) syscall_GetEvents(min_nr, nr int, events []IOEvent, timeout timespec) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}

	p := unsafe.Pointer(&events[0])
	n, _, err := syscall.Syscall6(
		syscall.SYS_IO_GETEVENTS,
		uintptr(c.ioctx),
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

func (c *AIOCtx) Write(p []byte) (n int, err error) {
	n, err = c.WriteAt(p, c.woff)
	if err != nil {
		return n, err
	}
	c.woff += int64(n)
	return n, nil
}

func (c *AIOCtx) WriteAt(p []byte, off int64) (n int, err error) {
	cb := newIOCB(c.fd)
	cb.prepareWrite(p, off)
	err = c.submit(cb)
	if err != nil {
		return 0, err
	}

	// 要等待

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
