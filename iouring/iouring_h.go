package iouring

type io_uring_sq struct {
	ring_sz  uint32
	ring_ptr uintptr
}

type io_uring_cq struct {
	ring_sz  uint32
	ring_ptr uintptr
}

// 这个玩意根据不同的setup参数会有不同的大小
/*
 * IO completion data structure (Completion Queue Entry)
 */
type io_uring_cqe struct {
	user_data uint64 /* sqe->data submission passed back */
	res       int32  /* result code for this event */
	flags     uint32

	// /*
	//  * If the ring is initialized with IORING_SETUP_CQE32, then this field
	//  * contains 16-bytes of padding, doubling the size of the CQE.
	//  */
	// __u64 big_cqe[]
}

// /*
//  * IO submission data structure (Submission Queue Entry)
//  */
//  struct io_uring_sqe {
// 	__u8	opcode;		/* type of operation for this sqe */
// 	__u8	flags;		/* IOSQE_ flags */
// 	__u16	ioprio;		/* ioprio for the request */
// 	__s32	fd;		/* file descriptor to do IO on */
// 	union {
// 		__u64	off;	/* offset into file */
// 		__u64	addr2;
// 		struct {
// 			__u32	cmd_op;
// 			__u32	__pad1;
// 		};
// 	};
// 	union {
// 		__u64	addr;	/* pointer to buffer or iovecs */
// 		__u64	splice_off_in;
// 		struct {
// 			__u32	level;
// 			__u32	optname;
// 		};
// 	};
// 	__u32	len;		/* buffer size or number of iovecs */
// 	union {
// 		__kernel_rwf_t	rw_flags;
// 		__u32		fsync_flags;
// 		__u16		poll_events;	/* compatibility */
// 		__u32		poll32_events;	/* word-reversed for BE */
// 		__u32		sync_range_flags;
// 		__u32		msg_flags;
// 		__u32		timeout_flags;
// 		__u32		accept_flags;
// 		__u32		cancel_flags;
// 		__u32		open_flags;
// 		__u32		statx_flags;
// 		__u32		fadvise_advice;
// 		__u32		splice_flags;
// 		__u32		rename_flags;
// 		__u32		unlink_flags;
// 		__u32		hardlink_flags;
// 		__u32		xattr_flags;
// 		__u32		msg_ring_flags;
// 		__u32		uring_cmd_flags;
// 		__u32		waitid_flags;
// 		__u32		futex_flags;
// 	};
// 	__u64	user_data;	/* data to be passed back at completion time */
// 	/* pack this to avoid bogus arm OABI complaints */
// 	union {
// 		/* index into fixed buffers, if used */
// 		__u16	buf_index;
// 		/* for grouped buffer selection */
// 		__u16	buf_group;
// 	} __attribute__((packed));
// 	/* personality to use, if used */
// 	__u16	personality;
// 	union {
// 		__s32	splice_fd_in;
// 		__u32	file_index;
// 		__u32	optlen;
// 		struct {
// 			__u16	addr_len;
// 			__u16	__pad3[1];
// 		};
// 	};
// 	union {
// 		struct {
// 			__u64	addr3;
// 			__u64	__pad2[1];
// 		};
// 		__u64	optval;
// 		/*
// 		 * If the ring is initialized with IORING_SETUP_SQE128, then
// 		 * this field is used for 80 bytes of arbitrary command data
// 		 */
// 		__u8	cmd[0];
// 	};
// };
