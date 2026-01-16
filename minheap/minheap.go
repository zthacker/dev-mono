package minheap

type MinHeap struct {
	data []int
}

func (h *MinHeap) Push(val int) {
	h.data = append(h.data, val)
	h.bubbleUp(len(h.data) - 1)
}

func (h *MinHeap) Pop() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}

	// grab min value first
	min := h.data[0]

	// move last to root
	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]

	// bubble down if there's still data
	if len(h.data) > 0 {
		h.bubbleDown(0)
	}

	return min, true
}

func (h *MinHeap) Peek() (int, bool) {
	if len(h.data) == 0 {
		return 0, false
	}

	return h.data[0], true
}

func (h *MinHeap) bubbleUp(idx int) {
	if idx == 0 {
		return
	}

	// val 0
	if h.data[idx] < h.data[h.parent(idx)] {
		// swap
		h.data[idx], h.data[h.parent(idx)] = h.data[h.parent(idx)], h.data[idx]
		h.bubbleUp(h.parent(idx))
	}
}

func (h *MinHeap) bubbleDown(idx int) {
	left := h.leftChild(idx)
	right := h.rightChild(idx)
	smallest := idx

	// check left
	if left < len(h.data) && h.data[left] < h.data[smallest] {
		smallest = left
	}

	// check right
	if right < len(h.data) && h.data[right] < h.data[smallest] {
		smallest = right
	}

	// if child is smaller, swap and continue
	if smallest != idx {
		h.data[idx], h.data[smallest] = h.data[smallest], h.data[idx]
		h.bubbleDown(smallest)
	}

}

func (h *MinHeap) parent(idx int) int     { return (idx - 1) / 2 }
func (h *MinHeap) leftChild(idx int) int  { return 2*idx + 1 }
func (h *MinHeap) rightChild(idx int) int { return 2*idx + 2 }
