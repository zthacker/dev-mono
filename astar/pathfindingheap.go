package astar

type NodeHeap struct {
	data []*Node
}

func (h *NodeHeap) Push(val *Node) {
	h.data = append(h.data, val)
	h.bubbleUp(len(h.data) - 1)
}

func (h *NodeHeap) Pop() (*Node, bool) {
	if len(h.data) == 0 {
		return nil, false
	}

	min := h.data[0]

	h.data[0] = h.data[len(h.data)-1]
	h.data = h.data[:len(h.data)-1]

	if len(h.data) > 0 {
		h.bubbleDown(0)
	}

	return min, true
}

func (h *NodeHeap) Len() int {
	return len(h.data)
}

func (h *NodeHeap) bubbleUp(idx int) {
	if idx == 0 {
		return
	}

	// val 0
	if h.data[idx].TotalEstimatedCost() < h.data[h.parent(idx)].TotalEstimatedCost() {
		// swap
		h.data[idx], h.data[h.parent(idx)] = h.data[h.parent(idx)], h.data[idx]
		h.bubbleUp(h.parent(idx))
	}
}

func (h *NodeHeap) bubbleDown(idx int) {
	left := h.leftChild(idx)
	right := h.rightChild(idx)
	smallest := idx

	// check left
	if left < len(h.data) && h.data[left].TotalEstimatedCost() < h.data[smallest].TotalEstimatedCost() {
		smallest = left
	}

	// check right
	if right < len(h.data) && h.data[right].TotalEstimatedCost() < h.data[smallest].TotalEstimatedCost() {
		smallest = right
	}

	// if child is smaller, swap and continue
	if smallest != idx {
		h.data[idx], h.data[smallest] = h.data[smallest], h.data[idx]
		h.bubbleDown(smallest)
	}

}

func (h *NodeHeap) parent(idx int) int     { return (idx - 1) / 2 }
func (h *NodeHeap) leftChild(idx int) int  { return 2*idx + 1 }
func (h *NodeHeap) rightChild(idx int) int { return 2*idx + 2 }
