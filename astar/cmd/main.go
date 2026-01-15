package main

import (
	"astar"
	"fmt"
)

func main() {
	//   0   1   2   3   4
	// ┌───┬───┬───┬───┬───┐
	// │ S │   │   │ # │   │ 0
	// ├───┼───┼───┼───┼───┤
	// │   │   │ # │ # │   │ 1
	// ├───┼───┼───┼───┼───┤
	// │   │   │   │   │   │ 2
	// ├───┼───┼───┼───┼───┤
	// │ # │ # │   │   │   │ 3
	// ├───┼───┼───┼───┼───┤
	// │   │   │   │   │ G │ 4
	// └───┴───┴───┴───┴───┘
	grid := &astar.Grid{
		Width:  5,
		Height: 5,
		Walls: map[[2]int]bool{
			{3, 0}: true,
			{2, 1}: true,
			{3, 1}: true,
			{1, 3}: true,
			{0, 3}: true,
		},
	}

	path := astar.FindPath(grid, 0, 0, 4, 4)

	if path == nil {
		fmt.Println("No path found")
		return
	}

	fmt.Println("Path found:")
	for _, node := range path {
		fmt.Printf("(%d, %d)\n", node.X, node.Y)
	}
}
