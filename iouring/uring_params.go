package iouring

const (
	IORING_SETUP_IOPOLL uint32 = 1 << iota
	IORING_SETUP_SQPOLL
	IORING_SETUP_SQ_AFF
	IORING_SETUP_CQSIZE
	IORING_SETUP_CLAMP
	IORING_SETUP_ATTACH_WQ
	IORING_SETUP_R_DISABLED
	IORING_SETUP_SUBMIT_ALL

	IORING_SETUP_COOP_TASKRUN

	IORING_SETUP_TASKRUN_FLAG
	IORING_SETUP_SQE128
	IORING_SETUP_CQE32
)

// Passed in for io_uring_setup(2). Copied back with updated info on success
type io_uring_params struct {
	sq_entries     uint32
	cq_entries     uint32
	flags          uint32
	sq_thread_cpu  uint32
	sq_thread_idle uint32
	features       uint32
	wq_fd          uint32
	resv           [3]uint32
	sq_off         io_sqring_offsets
	cq_off         io_cqring_offsets
}

type io_sqring_offsets struct {
	head         uint32
	tail         uint32
	ring_mask    uint32
	ring_entries uint32
	flags        uint32
	dropped      uint32
	array        uint32
	resv         [3]uint32
}

type io_cqring_offsets struct {
	head         uint32
	tail         uint32
	ring_mask    uint32
	ring_entries uint32
	overflow     uint32
	cqes         uint32
	flags        uint32
	resv         [3]uint32
}
