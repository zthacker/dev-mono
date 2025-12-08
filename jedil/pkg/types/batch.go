package types

import (
	"fmt"
)

// BATCH_SIZE is 4 to align with AVX2 (256-bit registers)
// 4 float64s * 64 bits = 256 bits
const BATCH_SIZE = 4

// Vec3Batch holds 4 vectors in Structure-of-Arrays (SoA) layout
// This is optimal for SIMD auto-vectorization
type Vec3Batch struct {
	Xs [BATCH_SIZE]float64
	Ys [BATCH_SIZE]float64
	Zs [BATCH_SIZE]float64
}

// NewVec3Batch creates a zeroed batch
func NewVec3Batch() Vec3Batch {
	return Vec3Batch{}
}

// Add performs component-wise addition on 4 vectors at once
func (v Vec3Batch) Add(other Vec3Batch) Vec3Batch {
	var result Vec3Batch
	for i := 0; i < BATCH_SIZE; i++ {
		result.Xs[i] = v.Xs[i] + other.Xs[i]
		result.Ys[i] = v.Ys[i] + other.Ys[i]
		result.Zs[i] = v.Zs[i] + other.Zs[i]
	}
	return result
}

// Sub performs component-wise subtraction
func (v Vec3Batch) Sub(other Vec3Batch) Vec3Batch {
	var result Vec3Batch
	for i := 0; i < BATCH_SIZE; i++ {
		result.Xs[i] = v.Xs[i] - other.Xs[i]
		result.Ys[i] = v.Ys[i] - other.Ys[i]
		result.Zs[i] = v.Zs[i] - other.Zs[i]
	}
	return result
}

// Dot calculates 4 dot products at once
// Returns a Batch of scalars (just one array of 4 floats)
func (v Vec3Batch) Dot(other Vec3Batch) [BATCH_SIZE]float64 {
	var result [BATCH_SIZE]float64
	for i := 0; i < BATCH_SIZE; i++ {
		result[i] = v.Xs[i]*other.Xs[i] +
			v.Ys[i]*other.Ys[i] +
			v.Zs[i]*other.Zs[i]
	}
	return result
}

// Scale multiplies all 4 vectors by 4 scalars
func (v Vec3Batch) Scale(scalars [BATCH_SIZE]float64) Vec3Batch {
	var result Vec3Batch
	for i := 0; i < BATCH_SIZE; i++ {
		result.Xs[i] = v.Xs[i] * scalars[i]
		result.Ys[i] = v.Ys[i] * scalars[i]
		result.Zs[i] = v.Zs[i] * scalars[i]
	}
	return result
}

// Helper to set a specific vector in the batch (for testing)
func (v *Vec3Batch) Set(index int, x, y, z float64) {
	if index >= 0 && index < BATCH_SIZE {
		v.Xs[index] = x
		v.Ys[index] = y
		v.Zs[index] = z
	}
}

// Helper to get a specific vector
func (v Vec3Batch) Get(index int) Vec3 {
	if index >= 0 && index < BATCH_SIZE {
		return Vec3{X: v.Xs[index], Y: v.Ys[index], Z: v.Zs[index]}
	}
	return Vec3{}
}

func (v Vec3Batch) String() string {
	return fmt.Sprintf("Batch[v0=%v, ...]", v.Get(0))
}
