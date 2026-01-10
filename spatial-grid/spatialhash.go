package spatialgrid

type CellKey struct{ X, Y int }

type SpatialHash struct {
	cellSize   float64
	cells      [][]Entity // 1D array
	gridWidth  int        // num of cells horizontally
	gridHeight int        // num of cells vertically

}

func NewSpatialHash(cellSize float64, worldWidth, worldHeight float64) *SpatialHash {

	gridHeight := int(worldHeight / cellSize)
	gridWidth := int(worldWidth / cellSize)

	return &SpatialHash{
		gridWidth:  gridHeight,
		gridHeight: gridWidth,
		cellSize:   cellSize,
		cells:      make([][]Entity, gridWidth+gridHeight),
	}
}

func (sh *SpatialHash) cellIndex(x, y float64) int {
	cellX := int(x / sh.cellSize)
	cellY := int(y / sh.cellSize)
	return cellX + cellY*sh.gridWidth
}

func (sh *SpatialHash) Insert(e Entity, x float64, y float64) {
	cellIndex := sh.cellIndex(x, y)
	sh.cells[cellIndex] = append(sh.cells[cellIndex], e)
}

func (sh *SpatialHash) Remove(e Entity) {
	cellIndex := sh.cellIndex(e.Position())
	entities := sh.cells[cellIndex]

	for i, entity := range entities {
		if entity.ID() == e.ID() {
			// swap since ordering doesn't matter in the cell
			entities[i] = entities[len(entities)-1]
			sh.cells[cellIndex] = entities[:len(entities)-1]
			return
		}
	}
}

func (sh *SpatialHash) Update(e Entity, x float64, y float64) {
	cellEntity := sh.cellIndex(e.Position())
	potentialNewCell := sh.cellIndex(x, y)
	if cellEntity != potentialNewCell {
		sh.Remove(e)
		sh.Insert(e, x, y)
	}

}

// All entities in cells that overlap the radius
func (sh *SpatialHash) Query(x float64, y float64, radius float64) []Entity {
	centerCellX := int(x / sh.cellSize)
	centerCellY := int(y / sh.cellSize)

	// how many  cells does this radius cover
	cellRadius := int(radius/sh.cellSize) + 1

	var entities []Entity
	for i := -cellRadius; i <= cellRadius; i++ {
		for j := -cellRadius; j <= cellRadius; j++ {
			cellX := centerCellX + i
			cellY := centerCellY + j

			// bounds check
			if cellX < 0 || cellX >= sh.gridWidth || cellY < 0 || cellY >= sh.gridHeight {
				continue
			}

			idx := cellX + cellY*sh.gridWidth
			entities = append(entities, sh.cells[idx]...)
		}
	}
	return entities
}
