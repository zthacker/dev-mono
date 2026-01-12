package spatialgrid

import "testing"

type TestEntity struct {
	id   int
	x, y float64
}

func (e *TestEntity) ID() int {
	return e.id
}

func (e *TestEntity) Position() (x, y float64) {
	return float64(e.x), float64(e.y)
}

func BenchmarkInsert(b *testing.B) {
	sh := NewSpatialHash(10, 1000, 1000)
	e := TestEntity{id: 1, x: 50, y: 50}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.Insert(&e, float64(e.x), float64(e.y))
	}
}

func BenchmarkQuery(b *testing.B) {
	sh := NewSpatialHash(10, 1000, 1000)

	for i := 0; i < 10000; i++ {
		x := float64(i%100) * 10
		y := float64(i/100) * 10
		sh.Insert(&TestEntity{id: i, x: x, y: y}, x, y)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.Query(500, 500, 50)
	}
}

func BenchmarkQueryDistanceFilter(b *testing.B) {
	sh := NewSpatialHash(10, 1000, 1000)

	for i := 0; i < 10000; i++ {
		x := float64(i%100) * 10
		y := float64(i/100) * 10
		sh.Insert(&TestEntity{id: i, x: x, y: y}, x, y)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sh.QueryDistanceFilter(500, 500, 50)
	}
}

func BenchmarkBruteForce(b *testing.B) {
	// All entities in one flat slice
	entities := make([]TestEntity, 10000)
	for i := 0; i < 10000; i++ {
		entities[i] = TestEntity{id: i, x: float64(i%100) * 10, y: float64(i/100) * 10}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var result []TestEntity
		for _, e := range entities {
			dx := e.x - 500
			dy := e.y - 500
			if dx*dx+dy*dy <= 50*50 {
				result = append(result, e)
			}
		}
	}
}
