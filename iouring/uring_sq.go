package iouring

type io_uring_sq struct {
	khead    *uint32
	ktail    *uint32
	kflags   *uint32
	kdropped *uint32
	array    *uint32
	sqes     io_uring_sqes

	sqe_head uint32
	sqe_tail uint32

	ring_sz  uint32
	ring_ptr uintptr

	ring_mask    *uint32
	ring_entries *uint32

	pad [2]uint32
}

type io_uring_sqes interface {
	entry_size() uint32
}

type io_uring_sqe_core struct {
	opcode      uint8
	flags       uint8
	ioprio      uint16
	fd          int32
	off         uint64
	addr        uint64
	len         uint32
	open_flags  uint32
	user_data   uint64
	buf_index   uint16
	personality uint16
	file_index  uint32
	optval      uint64
}

type io_uring_sqe_128 struct {
	io_uring_sqe_core
	cmd [80]uint8
}

type io_uring_sqe_64 struct {
	io_uring_sqe_core
	extra [2]uint64
}

type io_uring_sqes_128 struct {
	sqes []io_uring_sqe_128
}

func (e io_uring_sqes_128) entry_size() uint32 {
	return uint32(len(e.sqes)) * sizeof(io_uring_sqe_128{})
}

type io_uring_sqes_64 struct {
	sqes []io_uring_sqe_64
}

func (e io_uring_sqes_64) entry_size() uint32 {
	return uint32(len(e.sqes)) * sizeof(io_uring_sqes_64{})
}

func new_sqes(flags uint32) io_uring_sqes {
	if flags&IORING_SETUP_SQE128 == 0 {
		return io_uring_sqes_128{}
	} else {
		return io_uring_sqes_64{}
	}
}
