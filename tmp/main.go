package main

import (
	"fmt"
	"unsafe"
)

func main() {
	buf := []byte{1, 2, 3, 4, 5}
	ptr := &buf[0]

	ptr1 := uintptr(unsafe.Pointer(ptr))
	ptr2 := unsafe.Add(unsafe.Pointer(ptr1), 3)

	buf2 := (*byte)(ptr2)

	fmt.Println(*buf2)
}
