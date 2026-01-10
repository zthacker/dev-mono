package spatialgrid

type Entity interface {
	ID() int
	Position() (x, y float64)
}
