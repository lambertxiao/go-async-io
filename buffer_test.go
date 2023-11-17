package goaio

import (
	"errors"
	"testing"
	"unsafe"
)

func TestAlignment(t *testing.T) {
	buf, err := PosixMemAlign(1024*1024, 4096)
	if err != nil {
		t.Fatal(err)
	}

	if len(buf) != 1024*1024 {
		t.Fatal(errors.New("malloc buf size mismatch"))
	}

	if uintptr(unsafe.Pointer(&buf[0]))%4096 != 0 {
		t.Fatal(errors.New("mem align failed"))
	}
}
