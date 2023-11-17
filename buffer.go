package goaio

import (
	"errors"
	"unsafe"
)

const AlignSize = 4096

func PosixMemAlign(blockSize, alignSize uint) ([]byte, error) {
	// 判断blockSize是否是2的幂
	if alignSize != 0 && blockSize&(alignSize-1) != 0 {
		return nil, errors.New("invalid argument")
	}

	block := make([]byte, blockSize+alignSize)
	remainder := alignment(block, alignSize)
	var offset uint
	if remainder != 0 {
		offset = alignSize - remainder
	}
	return block[offset : offset+blockSize], nil
}

func alignment(block []byte, alignSize uint) uint {
	if len(block) < 1 {
		return 0
	}
	if alignSize == 0 || alignSize == 1 || alignSize&1 != 0 {
		return 0
	}
	return uint(uintptr(unsafe.Pointer(&block[0])) & uintptr(alignSize-1))
}
