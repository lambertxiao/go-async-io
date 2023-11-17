package libaio

import (
	goaio "go-aio"
	"os"
	"syscall"

	"github.com/hashicorp/go-multierror"
)

type AIOCtx struct {
	fd      *os.File
	ioctx   uint // 由内核libaio填充该iocxt
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

func OpenAIOCtx(fpath string, opts Options) (goaio.IOCtx, error) {
	fd, err := os.OpenFile(fpath, syscall.O_DIRECT|opts.Flag, opts.Perm)
	if err != nil {
		return nil, err
	}

	var ioctx uint
	err = syscall_iosetup(&ioctx, opts.IODepth)
	if err != nil {
		return nil, err
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
	n, err := syscall_getevents(c.ioctx, 1, 1, c.events, c.timeout)
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

	return 0, nil
}

func (c *AIOCtx) Read(p []byte) (n int, err error) {
	return 0, nil
}

func (c *AIOCtx) ReadAt(p []byte, off int64) (n int, err error) {
	return 0, nil
}

func (c *AIOCtx) Close() error {
	return syscall_destory(c.ioctx)
}
