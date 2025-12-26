package zone

import (
	"realm-server/internal/entity"
	"realm-server/pkg/math"
)

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
	cellSize float32          // how big each cell is -- i.e. 50 units
	cellsX   int              // how many cells across (columns)
	cellsZ   int              // how many cells deep (rows)
	bounds   math.WorldBounds // The world area this grid covers
	cells    []Cell           // all the cells stored in a flat array
}

// Cell contains entities within a spatial region.
type Cell struct {
	// Use a map for O(1) add/remove
	entities map[entity.EntityID]entity.Entity

	// Or use slices if entity count per cell is small
	// players []*entity.Player
	// npcs    []*entity.NPC
}

// CellCoord identifies a cell in the grid.
type CellCoord struct {
	X, Z int
}

func NewGrid(bounds math.WorldBounds, cellSize float32) *Grid {
	// Ex: World Bounds is Min(0,0,0) to Max(1000,100,600)
	// Cell Size: 50
	// Width (X) = 1000 - 0 = 1000 units
	// Depth (Z) = 600 - 0 = 600 units
	// Y is height, which isn't relevant in a 2D grid

	// Cell across (X)  1000 / 50 = 20 cells
	// Cells deep (Z) 600 / 50 = 12 cells
	// Total cells = 20 * 12 = 240 cells

	worldWidth := bounds.Max.X - bounds.Min.X
	worldDepth := bounds.Max.Z - bounds.Min.Z

	// int() truncates, so add 1 to handle partial cells
	cellsX := int(worldWidth/cellSize) + 1
	cellsZ := int(worldDepth/cellSize) + 1

	totalCells := cellsX * cellsZ
	cells := make([]Cell, totalCells)

	return &Grid{
		cellSize: cellSize,
		cellsX:   cellsX,
		cellsZ:   cellsZ,
		bounds:   bounds,
		cells:    cells,
	}
}

func (g *Grid) WorldToCell(pos math.Vec3) CellCoord {
	// - Convert world position to cell coordinates
	// - Clamp to valid range
}

func (g *Grid) CellToIndex(coord CellCoord) int {
	// - Convert 2D coord to 1D array index
	// - return coord.Z * g.cellsX + coord.X
}

func (g *Grid) AddEntity(e entity.Entity) {
	// - Get cell for entity position
	// - Add to cell's entity list
}

func (g *Grid) RemoveEntity(id entity.EntityID, pos math.Vec3) {
	// - Get cell for position
	// - Remove from cell
}

func (g *Grid) MoveEntity(e entity.Entity, oldPos, newPos math.Vec3) {
	// - Check if cell changed
	// - If so, remove from old, add to new
}

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
