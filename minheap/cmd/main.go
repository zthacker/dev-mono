package main

import (
	"fmt"
	"math/rand/v2"
	"minheap"
)

func main() {
	h := &minheap.MinHeap{}

	for i := 0; i < 100; i++ {
		h.Push(rand.IntN(100))
	}

	fmt.Println("Popping in order:")
	for {
		val, ok := h.Pop()
		if !ok {
			break
		}
		fmt.Println(val)
	}
}
