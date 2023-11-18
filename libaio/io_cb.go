package libaio

import (
	"os"
	"unsafe"
)

type IOCB struct {
	data      unsafe.Pointer
	key       uint64
	opcode    int16
	reqprio   int16
	fd        uint32
	buf       unsafe.Pointer
	nbytes    uint64
	offset    int64
	reserved2 int64
	flags     uint32
	resfd     uint32
}

func newIOCB(fd *os.File) *IOCB {
	return &IOCB{fd: uint32(fd.Fd()), reqprio: 0}
}

func (cb *IOCB) prepareWrite(buf []byte, offset int64) {
	if len(buf) <= 0 {
		return
	}

	p := unsafe.Pointer(&buf[0])
	cb.opcode = int16(IOCmdPwrite)
	cb.buf = p
	cb.nbytes = uint64(len(buf))
	cb.offset = offset
}

// func bytes2Iovec(bs [][]byte) []syscall.Iovec {
// 	var iovecs []syscall.Iovec
// 	for _, chunk := range bs {
// 		if len(chunk) == 0 {
// 			continue
// 		}
// 		iovecs = append(iovecs, syscall.Iovec{Base: &chunk[0]})
// 		iovecs[len(iovecs)-1].SetLen(len(chunk))
// 	}
// 	return iovecs
// }
