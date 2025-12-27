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

	// init each cell's map
	for i := range cells {
		cells[i].entities = make(map[entity.EntityID]entity.Entity)
	}

	return &Grid{
		cellSize: cellSize,
		cellsX:   cellsX,
		cellsZ:   cellsZ,
		bounds:   bounds,
		cells:    cells,
	}
}

func (g *Grid) WorldToCell(pos math.Vec3) CellCoord {
	// Offset position relative to grid origin
	relativeX := pos.X - g.bounds.Min.X
	relativeZ := pos.Z - g.bounds.Min.Z

	// Divide by cell size and truncate to get cell index
	cellX := int(relativeX / g.cellSize)
	cellZ := int(relativeZ / g.cellSize)

	// Clamp to valid range (handle positions outside bounds)
	if cellX < 0 {
		cellX = 0
	}
	if cellX >= g.cellsX {
		cellX = g.cellsX - 1
	}
	if cellZ < 0 {
		cellZ = 0
	}
	if cellZ >= g.cellsZ {
		cellZ = g.cellsZ - 1
	}

	return CellCoord{X: cellX, Z: cellZ}
}

// The formula:
// index = Z * cellsX + X
// Examples with cellsX = 4:
// Cell (0, 0): index = 0 * 4 + 0 = 0
// Cell (3, 0): index = 0 * 4 + 3 = 3
// Cell (0, 1): index = 1 * 4 + 0 = 4   <- Start of row 1
// Cell (2, 1): index = 1 * 4 + 2 = 6
// Cell (3, 2): index = 2 * 4 + 3 = 11
func (g *Grid) CellToIndex(coord CellCoord) int {
	// - Convert 2D coord to 1D array index
	return coord.Z*g.cellsX + coord.X
}

func (g *Grid) AddEntity(e entity.Entity) {
	// - Get cell for entity position
	entityCoord := g.WorldToCell(e.Movement().Position)

	// Index
	index := g.CellToIndex(entityCoord)

	// add to map
	g.cells[index].entities[e.ID()] = e

}

func (g *Grid) RemoveEntity(id entity.EntityID, pos math.Vec3) {
	// - Get cell for position
	entityCoord := g.WorldToCell(pos)

	// Index
	index := g.CellToIndex(entityCoord)

	// - Remove from cell
	delete(g.cells[index].entities, id)

}

func (g *Grid) MoveEntity(e entity.Entity, oldPos, newPos math.Vec3) {
	// - Check if cell changed
	oldCoord := g.WorldToCell(oldPos)
	newCoord := g.WorldToCell(newPos)

	// - If so, remove from old, add to new
	if oldCoord.X != newCoord.X || oldCoord.Z != newCoord.Z {
		// remove from old cell
		oldIndex := g.CellToIndex(oldCoord)
		delete(g.cells[oldIndex].entities, e.ID())

		// add to new cell
		newIndex := g.CellToIndex(newCoord)
		g.cells[newIndex].entities[e.ID()] = e
	}
}

// =============================================================================
// RANGE QUERIES
// =============================================================================

func (g *Grid) GetEntitiesInRange(center math.Vec3, radius float32) []entity.Entity {
	// get cells in range
	cellCoords := g.GetCellsInRange(center, radius)

	radiusSq := radius * radius

	var result []entity.Entity

	for _, coord := range cellCoords {
		index := g.CellToIndex(coord)
		cell := &g.cells[index]

		for _, e := range cell.entities {
			pos := e.Transform().Position

			// check distance
			distSq := center.DistanceSq(pos)
			if distSq <= radiusSq {
				result = append(result, e)
			}
		}
	}

	return result
}

func (g *Grid) GetCellsInRange(center math.Vec3, radius float32) []CellCoord {
	// Box corners
	minPos := math.Vec3{
		X: center.X - radius,
		Y: center.Y,
		Z: center.Z - radius,
	}

	maxPos := math.Vec3{
		X: center.X + radius,
		Y: center.Y,
		Z: center.Z + radius,
	}

	// convert to cell coords
	minCell := g.WorldToCell(minPos)
	maxCell := g.WorldToCell(maxPos)

	// collect all cells in the rectangle
	var cells []CellCoord
	for z := minCell.Z; z <= maxCell.Z; z++ {
		for x := minCell.X; x <= maxCell.X; x++ {
			cells = append(cells, CellCoord{X: x, Z: z})
		}
	}

	return cells
}

// =============================================================================
// NEIGHBOR ITERATION
// =============================================================================

func (g *Grid) GetNeighborCells(coord CellCoord, radius int) []CellCoord {
	var cells []CellCoord

	for z := coord.Z - radius; z <= coord.Z+radius; z++ {
		for x := coord.X - radius; x <= coord.X+radius; x++ {
			// Skip if out of bounds
			if x < 0 || x >= g.cellsX || z < 0 || z >= g.cellsZ {
				continue
			}
			cells = append(cells, CellCoord{X: x, Z: z})
		}
	}

	return cells
}

func (g *Grid) ForEachEntityInRadius(center math.Vec3, radius float32, fn func(entity.Entity)) {
	cellCoords := g.GetCellsInRange(center, radius)
	radiusSq := radius * radius

	for _, coord := range cellCoords {
		index := g.CellToIndex(coord)
		cell := &g.cells[index]

		for _, e := range cell.entities {
			pos := e.Transform().Position
			if center.DistanceSq(pos) <= radiusSq {
				fn(e)
			}
		}
	}
}
