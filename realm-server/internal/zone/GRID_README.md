# Spatial Grid System

## What is the Grid for?

Imagine your game world is a huge map (say 10,000 x 10,000 units). When a player moves, you need to answer: "Who is nearby?"

**Without a grid:** Check distance to every entity. 10,000 entities = 10,000 distance calculations. Every frame. For every player. Slow!

**With a grid:** Divide the world into cells. Only check entities in nearby cells. Maybe 50 entities instead of 10,000.

```
World map divided into cells:

+-------+-------+-------+-------+
|       |       |   @   |       |    @ = Player
|       |       |  ooo  |       |    o = Nearby entities (same cell)
+-------+-------+-------+-------+
|       |       |       |       |    x = Far entities (ignored)
|   x   |       |       |   x   |
+-------+-------+-------+-------+
|       |   x   |       |       |
|       |       |       |       |
+-------+-------+-------+-------+
```

---

## Grid Structure

```go
type Grid struct {
    cellSize float32        // How big each cell is (e.g., 50 units)
    cellsX   int            // How many cells across (columns)
    cellsZ   int            // How many cells deep (rows)
    bounds   math.WorldBounds // The world area this grid covers
    cells    []Cell         // All the cells stored in a flat array
}
```

---

## NewGrid - Creating the Grid

You need to figure out how many cells fit in the world bounds:

```
Example:
  World bounds: Min(0, 0, 0) to Max(1000, 100, 600)
  Cell size: 50 units

  World width (X):  1000 - 0 = 1000 units
  World depth (Z):  600 - 0 = 600 units
  (We ignore Y - it's height, not relevant for 2D grid)

  Cells across (X): 1000 / 50 = 20 cells
  Cells deep (Z):   600 / 50 = 12 cells

  Total cells: 20 * 12 = 240 cells
```

Visual:
```
         1000 units (X axis)
    <------------------------->
    +--+--+--+--+--+--+--+--+--+  ^
    |  |  |  |  |  |  |  |  |  |  |
    +--+--+--+--+--+--+--+--+--+  |
    |  |  |  |  |  |  |  |  |  |  | 600 units
    +--+--+--+--+--+--+--+--+--+  | (Z axis)
    |  |  |  |  |  |  |  |  |  |  |
    +--+--+--+--+--+--+--+--+--+  v

    Each cell is 50x50 units
```

---

## WorldToCell - Converting Position to Cell Coordinates

When an entity is at position `(327, 50, 183)`, which cell are they in?

```
World position: (327, 50, 183)
                  ^   ^    ^
                  X   Y    Z
                      |
                      ignored (height)
```

**The math:**

```
Cell size: 50 units

Cell X = (position.X - bounds.Min.X) / cellSize
       = (327 - 0) / 50
       = 6.54
       = 6  (truncate to int)

Cell Z = (position.Z - bounds.Min.Z) / cellSize
       = (183 - 0) / 50
       = 3.66
       = 3  (truncate to int)

Result: Cell (6, 3)
```

**Visual:**

```
        0     1     2     3     4     5     6     7
      +-----+-----+-----+-----+-----+-----+-----+-----+
   0  |     |     |     |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+-----+-----+
   1  |     |     |     |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+-----+-----+
   2  |     |     |     |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+-----+-----+
   3  |     |     |     |     |     |     |  @  |     |  <- Z=3
      +-----+-----+-----+-----+-----+-----+-----+-----+
                                            ^
                                           X=6

   @ = Entity at world position (327, 50, 183)
       is in cell (6, 3)
```

**Why subtract bounds.Min?**

If your world doesn't start at (0,0,0):

```
Bounds: Min(-500, 0, -500) to Max(500, 100, 500)
Position: (127, 50, 83)

Without offset: 127 / 50 = 2  (wrong!)

With offset: (127 - (-500)) / 50 = 627 / 50 = 12  (correct!)
```

The offset shifts coordinates so the grid starts at 0.

---

## CellToIndex - Converting 2D Coordinates to Flat Array Index

We store cells in a 1D array, but we think in 2D (X, Z). This function converts between them.

**The concept:**

```
2D Grid (how we think):        1D Array (how it's stored):

      X=0  X=1  X=2  X=3
    +----+----+----+----+       +----+----+----+----+----+----+----+----+----+----+----+----+
Z=0 |  0 |  1 |  2 |  3 |       |  0 |  1 |  2 |  3 |  4 |  5 |  6 |  7 |  8 |  9 | 10 | 11 |
    +----+----+----+----+       +----+----+----+----+----+----+----+----+----+----+----+----+
Z=1 |  4 |  5 |  6 |  7 |         ^                   ^                   ^
    +----+----+----+----+         |                   |                   |
Z=2 |  8 |  9 | 10 | 11 |       Row 0               Row 1               Row 2
    +----+----+----+----+       (Z=0)               (Z=1)               (Z=2)
```

Each row (Z value) is laid out sequentially. Row 0 takes indices 0-3, Row 1 takes 4-7, etc.

**The formula:**

```
index = Z * cellsX + X
```

**Examples with cellsX = 4:**

```
Cell (0, 0): index = 0 * 4 + 0 = 0
Cell (3, 0): index = 0 * 4 + 3 = 3
Cell (0, 1): index = 1 * 4 + 0 = 4   <- Start of row 1
Cell (2, 1): index = 1 * 4 + 2 = 6
Cell (3, 2): index = 2 * 4 + 3 = 11
```

**Visual breakdown for Cell (2, 1):**

```
      X=0  X=1  X=2  X=3
    +----+----+----+----+
Z=0 |    |    |    |    |   Skip Z=0 entirely (4 cells)
    +----+----+----+----+
Z=1 |    |    | @  |    |   Skip X=0, X=1 (2 cells)
    +----+----+----+----+
Z=2 |    |    |    |    |

Z * cellsX = 1 * 4 = 4  (skip 4 cells for row 0)
         + X = 2        (skip 2 cells in row 1)
             = 6        (final index)
```

**Why a flat array instead of 2D?**

```go
// 2D array (what you might expect):
cells[x][z]

// Flat array (what we use):
cells[z * cellsX + x]
```

Flat arrays are faster - better memory locality, simpler allocation.

---

## GetCellsInRange - Find which cells overlap a radius

Given a center point and radius, find all cells that could contain entities within that radius.

```
Center: (175, 0, 125)
Radius: 60 units
Cell size: 50 units

                    radius = 60
                  <----------->
        0     1     2     3     4     5
      +-----+-----+-----+-----+-----+-----+
   0  |     |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+
   1  |     |  +--+-----+--+  |     |     |
      +-----+--|  |     |  |--+-----+-----+   min corner
   2  |     |  |  |  @  |  |  |     |     |   (175-60, 125-60) = (115, 65)
      +-----+--|  |     |  |--+-----+-----+
   3  |     |  +--+-----+--+  |     |     |   max corner
      +-----+-----+-----+-----+-----+-----+   (175+60, 125+60) = (235, 185)
   4  |     |     |     |     |     |     |
      +-----+-----+-----+-----+-----+-----+

Cells in range: (2,1), (3,1), (2,2), (3,2), (2,3), (3,3)
```

**The algorithm:**

1. Calculate bounding box corners: center ± radius
2. Convert corners to cell coordinates
3. Return all cells in that rectangle

---

## GetEntitiesInRange - Find entities within radius

Uses `GetCellsInRange` and checks actual distances:

```
Same example, but now checking actual entities:

      +-----+-----+-----+-----+
      |     |  o  |     |     |    o = entity in cell but outside radius
      +-----+-----+-----+-----+
      |  o  |  x  |  @  |  x  |    x = entity in cell AND within radius
      +-----+-----+-----+-----+
      |     |  x  |     |  o  |    @ = center point
      +-----+-----+-----+-----+

We check all cells in the rectangle, but only return
entities whose actual distance is <= radius.
```

**Why DistanceSq instead of Distance?**

```go
// Slow - uses sqrt:
if center.Distance(pos) <= radius { ... }

// Fast - no sqrt:
if center.DistanceSq(pos) <= radius*radius { ... }
```

Since `sqrt(a) <= sqrt(b)` is the same as `a <= b`, we can compare squared values and skip the expensive square root. This matters when checking thousands of entities.

---

## GetNeighborCells - Get cells around a coordinate

Instead of world units, you specify a cell radius.

```
radius = 1 (3x3 grid of cells):

      +-----+-----+-----+
      | n   | n   | n   |     n = neighbor
      +-----+-----+-----+
      | n   | @   | n   |     @ = center cell
      +-----+-----+-----+
      | n   | n   | n   |
      +-----+-----+-----+

radius = 2 (5x5 grid of cells):

      +-----+-----+-----+-----+-----+
      | n   | n   | n   | n   | n   |
      +-----+-----+-----+-----+-----+
      | n   | n   | n   | n   | n   |
      +-----+-----+-----+-----+-----+
      | n   | n   | @   | n   | n   |
      +-----+-----+-----+-----+-----+
      | n   | n   | n   | n   | n   |
      +-----+-----+-----+-----+-----+
      | n   | n   | n   | n   | n   |
      +-----+-----+-----+-----+-----+
```

---

## ForEachEntityInRadius - Iterate without allocating

`GetEntitiesInRange` allocates a slice. If you're calling it every tick for every player, that's a lot of garbage collection. `ForEachEntityInRadius` uses a callback instead:

```go
// Allocating version - creates a slice
entities := grid.GetEntitiesInRange(playerPos, 100)
for _, e := range entities {
    sendUpdate(player, e)
}

// Non-allocating version - no slice created
grid.ForEachEntityInRadius(playerPos, 100, func(e entity.Entity) {
    sendUpdate(player, e)
})
```

Both do the same thing, but `ForEachEntityInRadius` produces less garbage. Use it in hot paths (like the tick loop).