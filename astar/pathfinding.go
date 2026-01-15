package astar

import "math"

// f(n) = g(n) + h(n)

// g(n) = actual cost from start to this node
// h(n) = estimated cost from this node to goal (heuristic)
// f(n) = total estimated cost (lower = better)

type Node struct {
	X, Y      int     // position
	Cost      float64 // cost from start
	Heuristic float64 // estimated cost to goal
	Parent    *Node   // reconstruct path
}

func (n *Node) TotalEstimatedCost() float64 {
	return n.Cost + n.Heuristic
}

type Grid struct {
	Width, Height int
	Walls         map[[2]int]bool // blocked cells
}

func manhattanDistance(x1, y1, x2, y2 int) float64 {
	return math.Abs(float64(x2-x1)) + math.Abs(float64(y2-y1))
}

func euclideanDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x2 - x1)
	dy := float64(y2 - y1)
	return math.Sqrt(dx*dx + dy*dy)
}

func getNeighborsFourDirection(grid *Grid, node *Node) [][2]int {
	directions := [][2]int{
		{0, -1}, // up
		{0, 1},  // down
		{-1, 0}, // left
		{1, 0},  // right
	}

	var neighbors [][2]int

	// check all 4 directions
	for _, dir := range directions {
		x := node.X + dir[0]
		y := node.Y + dir[1]

		// check bounds
		if x < 0 || x >= grid.Width || y < 0 || y >= grid.Height {
			continue
		}

		// wall check
		if grid.Walls[[2]int{x, y}] {
			continue
		}

		neighbors = append(neighbors, [2]int{x, y})
	}

	return neighbors
}

// Goal → Parent → Parent → Parent → Start (Parent = nil)

// Then reverse it to get:

// Start → ... → ... → Goal
func reconstructPath(node *Node) []*Node {
	var path []*Node

	// walk backwards from goal to start
	for node != nil {
		path = append(path, node)
		node = node.Parent
	}

	// reverse to get start -> goal order
	// could also just use slice.Reverse(path), but this is fun
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}

	return path
}

func FindPath(grid *Grid, startX, startY, goalX, goalY int) []*Node {
	// openList are the nodes to explore
	openList := []*Node{}

	// closed set are the nodes already explores
	closedSet := make(map[[2]int]bool)

	// look up for nodes by position to update them
	nodeMap := make(map[[2]int]*Node)

	// create the start node
	startNode := &Node{
		X:         startX,
		Y:         startY,
		Cost:      0,
		Heuristic: manhattanDistance(startX, startY, goalX, goalY),
		Parent:    nil,
	}

	openList = append(openList, startNode)
	nodeMap[[2]int{startX, startY}] = startNode

	// start our loop
	for len(openList) > 0 {
		// find nod with lowest cost
		lowestIdx := 0
		for i, n := range openList {
			if n.TotalEstimatedCost() < openList[lowestIdx].TotalEstimatedCost() {
				lowestIdx = i
			}
		}
		current := openList[lowestIdx]

		// see if it's the goal
		if current.X == goalX && current.Y == goalY {
			return reconstructPath(current)
		}

		// remove from open and add to close
		openList[lowestIdx] = openList[len(openList)-1]
		openList = openList[:len(openList)-1]

		// add to closed set
		closedSet[[2]int{current.X, current.Y}] = true

		for _, pos := range getNeighborsFourDirection(grid, current) {
			// skip
			if closedSet[pos] {
				continue
			}

			possibleCost := current.Cost + 1 // cost to reach neighbor

			// is this in openlist?
			existing := nodeMap[pos]
			if existing == nil {
				// New node for open list
				neighbor := &Node{
					X:         pos[0],
					Y:         pos[1],
					Cost:      possibleCost,
					Heuristic: manhattanDistance(pos[0], pos[1], goalX, goalY),
					Parent:    current,
				}
				openList = append(openList, neighbor)
				nodeMap[pos] = neighbor
			} else if possibleCost < existing.Cost {
				// found a better path
				existing.Cost = possibleCost
				existing.Parent = current
			}

		}
	}

	return nil // no path found -- womp
}
