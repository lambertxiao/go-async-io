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
	params  *io_uring_params
	sq      *io_uring_sq
	cq      *io_uring_cq
}

func New(entries uint) (goaio.IOCtx, error) {
	c := IOUringCtx{
		params: &io_uring_params{},
	}
	fd, err := syscall_io_uring_setup(entries, c.params)
	if err != nil {
		return nil, err
	}

	c.ring_fd = fd
	err = c.io_uring_mmap()
	if err != nil {
		return nil, err
	}
	return nil, nil
}

// io_uring 通过用户态与内核态共享内存的方式，来免去了使用系统调用发起 I/O 操作的过程
// 需要分别为提交队列SQ, 完成队列CQ, 提交队列项数组SQE进行mmap
func (c *IOUringCtx) io_uring_mmap() error {
	if err := c.mmap_sq(); err != nil {
		return err
	}
	if err := c.mmap_cq(); err != nil {
		return err
	}

	return nil
}

// 创建SQ内存映射
func (c *IOUringCtx) mmap_sq() error {
	c.sq = new(io_uring_sq)
	sq := c.sq
	sq.ring_sz = c.params.sq_off.array + c.params.sq_entries*uint32(unsafe.Sizeof(uint32(0)))
	ptr, err := syscall_mmap(c.ring_fd, sq.ring_sz, IORING_OFF_SQ_RING)
	if err != nil {
		return err
	}
	sq.ring_ptr = ptr

	// 将sq的属性映射到对应的内存区域
	sq.khead = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.head)
	sq.ktail = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.tail)
	sq.ring_mask = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.ring_mask)
	sq.ring_entries = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.ring_entries)
	sq.kflags = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.flags)
	sq.kdropped = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.dropped)
	sq.array = ptr_add_uint32(sq.ring_ptr, c.params.sq_off.array)

	seqs := new_sqes(c.params.flags)
	ptr, err = syscall_mmap(c.ring_fd, seqs.entry_size()*c.params.sq_entries, IORING_OFF_SQES)
	if err != nil {
		return err
	}

	c.sq.sqes = *(*io_uring_sqes)(ptr_to_pointer(ptr))
	return nil
}

// 创建CQ内存映射
func (c *IOUringCtx) mmap_cq() error {
	c.cq = new(io_uring_cq)
	cq := c.cq

	cqes := new_cqes(c.params.flags)
	cq.ring_sz = c.params.cq_off.cqes + c.params.cq_entries*cqes.entry_size()

	ptr, err := syscall_mmap(c.ring_fd, cq.ring_sz, IORING_OFF_CQ_RING)
	if err != nil {
		return err
	}
	cq.ring_ptr = ptr

	cq.khead = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.head)
	cq.ktail = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.tail)
	cq.ring_mask = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.ring_mask)
	cq.ring_entries = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.ring_entries)
	cq.kflags = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.flags)
	cq.koverflow = ptr_add_uint32(cq.ring_ptr, c.params.cq_off.overflow)
	cq.cqes = *(*io_uring_cqes)(ptr_add_pointer(cq.ring_ptr, c.params.cq_off.cqes))

	return nil
}

func (c *IOUringCtx) loop() {}
