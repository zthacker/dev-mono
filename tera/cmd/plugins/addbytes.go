package main

import (
	"unsafe"
)

//export allocate_buffer
func allocate_buffer(size uint32) uint32 {
	buf := make([]byte, size)
	ptr := uint32(uintptr(unsafe.Pointer(&buf[0])))

	return ptr
}

//export process_packet
func process_packet(ptr uint32, size uint32) uint64 {
	// Read the input data
	data := unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size)

	// Process in place; or write back to shared buffer
	// Copy to avoid aliasing
	result := append([]byte(nil), data...)

	// Do our byte adding
	result = append(result, 0xEE, 0xFF)

	// Write result back to shared buffer (if it fits)
	if len(result) <= int(size) {
		copy(unsafe.Slice((*byte)(unsafe.Pointer(uintptr(ptr))), size), result)
		return (uint64(ptr<<32) | uint64(len(result)))
	}

	// If the result is larger, use a temp buffer
	// The host (caller to this VM) reads it immediately, so GC here is ok
	newPtr := uint32(uintptr(unsafe.Pointer(&result[0])))
	return (uint64(newPtr<<32) | uint64(len(result)))
}

// Required for TinyGo to compile (even though we don't use it)
func main() {}
