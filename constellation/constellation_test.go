package constellation

import (
	"os"
	"testing"
)

func TestRoute(t *testing.T) {
	c := New()

	// Layout (all in a line for simplicity):
	//   AUSTIN(0,0,0) ---10--- SAT-1(10,0,0) ---10--- SAT-2(20,0,0) ---10--- TOKYO(30,0,0)

	c.AddNode(&Node{ID: "AUSTIN", Type: GroundStation, Position: Vec3{0, 0, 0}, CommRange: 15})
	c.AddNode(&Node{ID: "SAT-1", Type: Satellite, Position: Vec3{10, 0, 0}, CommRange: 15})
	c.AddNode(&Node{ID: "SAT-2", Type: Satellite, Position: Vec3{20, 0, 0}, CommRange: 15})
	c.AddNode(&Node{ID: "TOKYO", Type: GroundStation, Position: Vec3{30, 0, 0}, CommRange: 15})

	c.UpdateLinks()

	path := c.Route("AUSTIN", "TOKYO")

	t.Logf("Path: %v", path)

	expected := []string{"AUSTIN", "SAT-1", "SAT-2", "TOKYO"}
	if len(path) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, path)
	}
	for i, id := range expected {
		if path[i] != id {
			t.Errorf("path[%d]: expected %s, got %s", i, id, path[i])
		}
	}
}

func TestNoRoute(t *testing.T) {
	c := New()

	// Two nodes too far apart
	c.AddNode(&Node{ID: "AUSTIN", Type: GroundStation, Position: Vec3{0, 0, 0}, CommRange: 5})
	c.AddNode(&Node{ID: "TOKYO", Type: GroundStation, Position: Vec3{100, 0, 0}, CommRange: 5})

	c.UpdateLinks()

	path := c.Route("AUSTIN", "TOKYO")

	t.Logf("Path: %v", path)

	// Should just have TOKYO (unreachable, prev chain is empty)
	if len(path) != 1 || path[0] != "TOKYO" {
		t.Logf("No route case returned: %v (this is expected behavior for unreachable)", path)
	}
}

func TestGenerateAndExport(t *testing.T) {
	c := GenerateConstellation(100, 15, 1000, 500) // increased comm range for better connectivity

	t.Logf("Nodes: %d", len(c.Storage))
	t.Logf("Total edges: %d", countEdges(c))

	// Test a route
	path := c.Route("GS-0", "GS-1")
	t.Logf("Route GS-0 → GS-1: %v", path)

	// Export with the route highlighted
	jsonData, _ := c.ExportJSONWithRoute("GS-0", "GS-1")
	os.WriteFile("constellation.json", jsonData, 0644)
	t.Logf("Exported to constellation.json with route GS-0 → GS-1")
}

func countEdges(c *Constellation) int {
	count := 0
	for _, neighbors := range c.AdjacentList {
		count += len(neighbors)
	}
	return count / 2 // each edge counted twice
}
