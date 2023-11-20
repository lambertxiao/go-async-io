package iouring

import "unsafe"

func ptr_to_pointer(p uintptr) unsafe.Pointer {
	return unsafe.Pointer(p)
}

func ptr_add_pointer(p uintptr, len uint32) unsafe.Pointer {
	return unsafe.Add(ptr_to_pointer(p), len)
}

func ptr_add_uint32(p uintptr, len uint32) *uint32 {
	return (*uint32)(ptr_add_pointer(p, len))
}

func sizeof(obj interface{}) uint32 {
	return (uint32)(unsafe.Sizeof(obj))
}
