package iouring

import (
	goaio "go-async-io"
	"unsafe"
)

const (
	IORING_OFF_SQ_RING uint64 = 0
	IORING_OFF_CQ_RING uint64 = 0x8000000
	IORING_OFF_SQES    uint64 = 0x10000000
)

type IOUringCtx struct {
	ring_fd int
	params  *IOUringParams
	sq      *io_uring_sq
	cq      *io_uring_cq
}

func New(entries uint) (goaio.IOCtx, error) {
	c := IOUringCtx{
		params: &IOUringParams{},
	}
	fd, err := syscall_io_uring_setup(entries, c.params)
	if err != nil {
		return nil, err
	}

	c.ring_fd = fd
	c.mmapIOUring()
	go c.loop()
	return nil, nil
}

// io_uring 通过用户态与内核态共享内存的方式，来免去了使用系统调用发起 I/O 操作的过程
// io_uring 主要创建了 3 块共享内存, 提交队列SQ, 完成队列CQ, 提交队列项数组SQE
func (c *IOUringCtx) mmapIOUring() error {
	c.sq = new(io_uring_sq)
	c.cq = new(io_uring_cq)

	sq := c.sq
	cq := c.cq

	sq.ring_sz = c.params.sq_off.array + c.params.sq_entries*uint32(unsafe.Sizeof(uint32(0)))
	cq.ring_sz = c.params.cq_off.cqes + c.params.cq_entries*uint32(unsafe.Sizeof(io_uring_cqe{}))

	ptr, err := syscall_mmap(c.ring_fd, sq.ring_sz, IORING_OFF_SQ_RING)
	if err != nil {
		return err
	}
	sq.ring_ptr = ptr

	ptr, err = syscall_mmap(c.ring_fd, cq.ring_sz, IORING_OFF_CQ_RING)
	if err != nil {
		return err
	}
	cq.ring_ptr = ptr

	size := uint32(unsafe.Sizeof(io_uring_sqe{}))
	ptr, err = syscall_mmap(c.ring_fd, size*c.params.sq_entries, IORING_OFF_SQES)
	if err != nil {
		return err
	}

	c.sq.sqes = *(*[]io_uring_sqe)(unsafe.Pointer(ptr))

	return nil
}

// 创建SQ内存映射
func (c *IOUringCtx) mmapSQ() error {
	return nil
}

// 创建CQ内存映射
func (c *IOUringCtx) mmapCQ() error {
	return nil
}

// 创建SQE内存映射
func (c *IOUringCtx) mmapSQE() error {
	return nil
}

func (c *IOUringCtx) loop() {}
