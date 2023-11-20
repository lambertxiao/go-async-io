package iouring

type io_uring_cq struct {
	khead        *uint32
	ktail        *uint32
	kflags       *uint32
	koverflow    *uint32
	cqes         io_uring_cqes
	ring_sz      uint32
	ring_ptr     uintptr
	ring_mask    *uint32
	ring_entries *uint32
	pad          [2]uint32
}

type io_uring_cqes interface {
	entry_size() uint32
}

type io_uring_cqe_core struct {
	user_data uint64
	res       int32
	flags     uint32
}

type io_uring_cqe_16 struct {
	io_uring_cqe_core
}

type io_uring_cqes_16 struct {
	cqes []io_uring_cqe_16
}

func (q io_uring_cqes_16) entry_size() uint32 {
	return sizeof(io_uring_cqe_16{})
}

type io_uring_cqe_32 struct {
	io_uring_cqe_core
	extra1 uint64
	extra2 uint64
}

type io_uring_cqes_32 struct {
	cqes []io_uring_cqe_32
}

func (q io_uring_cqes_32) entry_size() uint32 {
	return sizeof(io_uring_cqe_32{})
}

func new_cqes(flags uint32) io_uring_cqes {
	if flags&IORING_SETUP_CQE32 == 0 {
		return io_uring_cqes_32{}
	} else {
		return io_uring_cqes_16{}
	}
}
