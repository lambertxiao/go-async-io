package libaio

import (
	"errors"
	goaio "go-async-io"
	"log"
	"os"
	"sync"
	"syscall"
	"unsafe"

	"github.com/hashicorp/go-multierror"
)

type AsyncIOCtx struct {
	fd      *os.File
	ioctx   uint // 由内核libaio填充该iocxt
	woff    int64
	roff    int64
	events  []IOEvent
	timeout timespec

	activeLock sync.RWMutex
	activeIOs  map[unsafe.Pointer]*ActiveIO
	closeCh    chan struct{}
}

type Options struct {
	IODepth int
	Flag    int
	Perm    os.FileMode
	// libaio timeout, unit ms 0 means no timeout.
	Timeout int
}

func New(fpath string, opts Options) (goaio.IOCtx, error) {
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

	ctx := &AsyncIOCtx{
		fd:        fd,
		ioctx:     ioctx,
		events:    events,
		timeout:   t,
		activeIOs: make(map[unsafe.Pointer]*ActiveIO),
		closeCh:   make(chan struct{}),
	}

	go ctx.loop()
	return ctx, nil
}

func (c *AsyncIOCtx) makeActiveIO(cb *IOCB) *ActiveIO {
	acio := newActiveIO(cb)
	c.activeLock.Lock()
	defer c.activeLock.Unlock()

	c.activeIOs[unsafe.Pointer(cb)] = acio
	return acio
}

func (c *AsyncIOCtx) removeActiveIO(cb *IOCB) {
	c.activeLock.Lock()
	defer c.activeLock.Unlock()

	delete(c.activeIOs, unsafe.Pointer(cb))
}

func (c *AsyncIOCtx) loop() {
	for {
		select {
		case <-c.closeCh:
			c.waitEvents()
			c.closeCh <- struct{}{}
			return
		default:
			c.waitEvents()
		}
	}
}

func (c *AsyncIOCtx) waitEvents() error {
	n, err := syscall_getevents(c.ioctx, 1, 2, c.events, c.timeout)
	if err != nil {
		return err
	}

	var errs error
	for i := 0; i < n; i++ {
		err := c.parseDoneEvent(c.events[i])
		if err != nil {
			errs = multierror.Append(errs, err)
		}
	}
	return errs
}

func (c *AsyncIOCtx) parseDoneEvent(event IOEvent) error {
	c.activeLock.RLock()
	defer c.activeLock.RUnlock()

	log.Println("parse event")
	io, ok := c.activeIOs[unsafe.Pointer(event.cb)]
	if !ok {
		return errors.New("event cb is not found")
	}

	io.retBytes = event.res
	io.Done()
	return nil
}

func (c *AsyncIOCtx) submitIO(cmd IOCmd, data []byte, off int64) (n int, err error) {
	cb := newIOCB(c.fd)

	switch cmd {
	case IOCmdPwrite:
		cb.prepareWrite(data, off)
	case IOCmdPread:
		cb.prepareRead(data, off)
	default:
		return 0, errors.New("unsupport cmd")
	}

	acio := c.makeActiveIO(cb)
	err = syscall_submit(c.ioctx, cb)
	if err != nil {
		c.removeActiveIO(cb)
		return 0, err
	}

	acio.Wait()
	return int(acio.retBytes), nil
}

func (c *AsyncIOCtx) Write(p []byte) (n int, err error) {
	n, err = c.WriteAt(p, c.woff)
	if err != nil {
		return n, err
	}
	c.woff += int64(n)
	return n, nil
}

func (c *AsyncIOCtx) WriteAt(p []byte, off int64) (n int, err error) {
	return c.submitIO(IOCmdPwrite, p, off)
}

func (c *AsyncIOCtx) Read(p []byte) (n int, err error) {
	n, err = c.ReadAt(p, c.woff)
	if err != nil {
		return n, err
	}
	c.roff += int64(n)
	return n, nil
}

func (c *AsyncIOCtx) ReadAt(p []byte, off int64) (n int, err error) {
	return c.submitIO(IOCmdPread, p, off)
}

func (c *AsyncIOCtx) Close() error {
	c.closeCh <- struct{}{}
	<-c.closeCh

	return syscall_destory(c.ioctx)
}
