package main

import (
	"io"
	lockfreering "lockfree-ring"
	"sync"
)

func main() {
	ring := lockfreering.NewRing(5)

	var wg sync.WaitGroup

	wg.Add(1)
	go func(r *lockfreering.Ring) {
		defer wg.Done()
		for i := range 10000 {
			for !ring.Push(i) {
				// buffer full so spin until consumer catches up
			}
		}

	}(ring)

	wg.Add(1)
	go func(r *lockfreering.Ring) {
		defer wg.Done()
		for i := 0; i < 10000; i++ {
			for {
				if val, ok := ring.Pop(); ok {
					println(val.(int))
					break
				}
				// buffer empty so spin until producer pushes something
			}
		}

	}(ring)
	wg.Wait()

	r := io.Reader()
}
