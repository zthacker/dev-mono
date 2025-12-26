package zone

// Uncomment when implementing:
// import (
// 	"realm-server/internal/entity"
// 	"realm-server/pkg/math"
// )

// Grid implements spatial partitioning for efficient neighbor queries.
// The zone is divided into cells; each cell tracks which entities are in it.
//
// Why spatial partitioning matters:
// - Without it: "Who's near me?" = O(n) check against all entities
// - With grid: O(1) cell lookup + check only nearby cells
// - Critical for AOI (Area of Interest) - what entities a player can see
//
// Cell size tradeoffs:
// - Small cells (10-20 units): Precise queries, more cells to check for large ranges
// - Large cells (100+ units): Fewer cells, but more entities per cell
// - Typical: 30-50 units (matches view distance / 3)
type Grid struct {
	cellSize float32
	cellsX   int
	cellsZ   int
	// bounds math.WorldBounds // Uncomment when implementing
	cells []Cell
}

// Cell contains entities within a spatial region.
type Cell struct {
	// Use a map for O(1) add/remove
	// entities map[entity.EntityID]entity.Entity

	// Or use slices if entity count per cell is small
	// players []*entity.Player
	// npcs    []*entity.NPC
}

// CellCoord identifies a cell in the grid.
type CellCoord struct {
	X, Z int
}

// TODO: Implement Grid:
//
// func NewGrid(bounds math.WorldBounds, cellSize float32) *Grid
//   - Calculate cellsX, cellsZ from bounds and cellSize
//   - Allocate cells slice
//
// func (g *Grid) WorldToCell(pos math.Vec3) CellCoord
//   - Convert world position to cell coordinates
//   - Clamp to valid range
//
// func (g *Grid) CellToIndex(coord CellCoord) int
//   - Convert 2D coord to 1D array index
//   - return coord.Z * g.cellsX + coord.X
//
// func (g *Grid) AddEntity(e entity.Entity)
//   - Get cell for entity position
//   - Add to cell's entity list
//
// func (g *Grid) RemoveEntity(id entity.EntityID, pos math.Vec3)
//   - Get cell for position
//   - Remove from cell
//
// func (g *Grid) MoveEntity(e entity.Entity, oldPos, newPos math.Vec3)
//   - Check if cell changed
//   - If so, remove from old, add to new

// =============================================================================
// RANGE QUERIES
// =============================================================================

// TODO: Implement range queries:
//
// func (g *Grid) GetEntitiesInRange(center math.Vec3, radius float32) []entity.Entity
//   Algorithm:
//   1. Calculate cell range: (center-radius) to (center+radius)
//   2. For each cell in range:
//      - For each entity in cell:
//        - If distance <= radius, add to result
//   3. Return results
//
//   Optimization: Use radiusSq and DistanceSq to avoid sqrt
//
// func (g *Grid) GetCellsInRange(center math.Vec3, radius float32) []CellCoord
//   - Just return the cell coordinates, let caller iterate

// =============================================================================
// NEIGHBOR ITERATION
// =============================================================================

// For AOI updates, you often need to iterate neighboring cells.

// TODO: Implement neighbor helpers:
//
// func (g *Grid) GetNeighborCells(coord CellCoord, radius int) []CellCoord
//   - Return cells within 'radius' cells of coord
//   - radius=1 gives 9 cells (3x3)
//   - radius=2 gives 25 cells (5x5)
//
// func (g *Grid) ForEachEntityInRadius(center math.Vec3, radius float32, fn func(entity.Entity))
//   - Efficient iteration without allocating result slice
