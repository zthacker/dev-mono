package lockfreering

import (
	"log"
	"sync/atomic"
)

type Ring struct {
	buffer []any
	head   atomic.Uint64 // pad to 64 bytes so we can have a separate cache line for head
	_      [56]byte
	tail   atomic.Uint64 // pad to 64 bytes so we can have a separate cache line for tail
	_      [56]byte
	size   int
}

func NewRing(size int) *Ring {
	return &Ring{
		buffer: make([]any, size),
		head:   atomic.Uint64{},
		tail:   atomic.Uint64{},
		size:   size,
	}
}

func (r *Ring) Push(val any) bool {
	if (int(r.tail.Load())+1)%r.size == int(r.head.Load()) {
		log.Print("full")
		return false
	}

	r.buffer[r.tail.Load()] = val
	r.tail.Store((r.tail.Load() + 1) % uint64(r.size))

	return true
}

func (r *Ring) Pop() (any, bool) {
	if r.head.Load() == r.tail.Load() {
		log.Print("empty")
		return nil, false
	}
	val := r.buffer[r.head.Load()]
	r.head.Store((r.head.Load() + 1) % uint64(r.size))
	return val, true
}
