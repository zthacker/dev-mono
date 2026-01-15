package astar

import "testing"

func BenchmarkFindPath_SmallGrid(b *testing.B) {
	grid := &Grid{
		Width:  10,
		Height: 10,
		Walls:  map[[2]int]bool{{3, 0}: true, {3, 1}: true, {3, 2}: true},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindPath(grid, 0, 0, 9, 9)
	}
}

func BenchmarkFindPath_LargeGrid(b *testing.B) {
	grid := &Grid{
		Width:  100,
		Height: 100,
		Walls:  make(map[[2]int]bool),
	}
	// Add some walls
	for i := 0; i < 50; i++ {
		grid.Walls[[2]int{50, i}] = true
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindPath(grid, 0, 0, 99, 99)
	}
}

func BenchmarkFindPath_MazeGrid(b *testing.B) {
	grid := &Grid{
		Width:  50,
		Height: 50,
		Walls:  make(map[[2]int]bool),
	}
	// Create a maze-like pattern
	for y := 0; y < 50; y += 2 {
		for x := 0; x < 48; x++ {
			grid.Walls[[2]int{x, y}] = true
		}
	}
	for y := 1; y < 50; y += 2 {
		for x := 2; x < 50; x++ {
			grid.Walls[[2]int{x, y}] = true
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FindPath(grid, 0, 1, 49, 49)
	}
}
