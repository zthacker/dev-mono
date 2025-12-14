package main

import (
	"unsafe"
)

//export allocate_buffer
func allocate_buffer(size uint32) uint32 {
	buf := make([]byte, size)
	return uint32(uintptr(unsafe.Pointer(&buf[0])))
}

//export process_packet
func process_packet(ptr uint32, size uint32) uint64 {
	data := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)

	data = append(data, 0xEE, 0xFF)

	newPtr := uint32(uintptr(unsafe.Pointer(&data[0])))
	newLen := uint32(len(data))

	return (uint64(newPtr) << 32) | uint64(newLen)
}

// Required for TinyGo to compile (even though we don't use it)
func main() {}
