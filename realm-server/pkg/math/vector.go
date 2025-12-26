package math

import "math"

// Vec3 represents a 3D position or direction in world space.
// Uses float32 for network efficiency (matches most game clients).
type Vec3 struct {
	X, Y, Z float32
}

// - Add(other Vec3) Vec3
func (v *Vec3) Add(vec3 Vec3) {
	v.X += vec3.X
	v.Y += vec3.Y
	v.Z += vec3.Z
}

// - Sub(other Vec3) Vec3
func (v *Vec3) Sub(vec3 Vec3) {
	v.X -= vec3.X
	v.Y -= vec3.Y
	v.Z -= vec3.Z
}

// - Scale(s float32) Vec3
func (v *Vec3) Scale(s float32) {
	v.X *= s
	v.Y *= s
	v.Z *= s
}

// - Length() float32
func (v *Vec3) Length() float32 {
	numsSquared := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return float32(math.Sqrt(float64(numsSquared)))
}

// - LengthSq() float32 - useful for distance comparisons without sqrt
func (v *Vec3) LengthSq() float32 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

// - Normalize() Vec3
func (v *Vec3) Normalize() {
	l := v.Length()

	if l == 0 {
		return
	}

	v.X /= l
	v.Y /= l
	v.Z /= l

}

// - Distance(other Vec3) float32
func (v *Vec3) Distance(other Vec3) float32 {
	return float32(math.Sqrt(float64(v.DistanceSq(other))))
}

// - DistanceSq(other Vec3) float32
func (v *Vec3) DistanceSq(other Vec3) float32 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	dz := v.Z - other.Z
	return dx*dx + dy*dy + dz*dz
}

// - Distance2D(other Vec3) float32 - ignores Y for ground distance
func (v *Vec3) Distance2D(other Vec3) float32 {
	dx := v.X - other.X
	dz := v.Z - other.Z
	return float32(math.Sqrt(float64(dx*dx + dz*dz)))
}

// - Dot(other Vec3) float32
func (v *Vec3) Dot(other Vec3) float32 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// WorldBounds represents an axis-aligned bounding box.
// Used for zone boundaries, grid cells, and spatial queries.
type WorldBounds struct {
	Min, Max Vec3
}

// - Contains(p Vec3) bool
func (w *WorldBounds) Contains(p Vec3) bool {
	return p.X >= w.Min.X && p.X <= w.Max.X &&
		p.Y >= w.Min.Y && p.Y <= w.Max.Y &&
		p.Z >= w.Min.Z && p.Z <= w.Max.Z

}

// - Intersects(other WorldBounds) bool
func (w *WorldBounds) Intersects(other WorldBounds) bool {
	return (w.Min.X <= other.Max.X && w.Max.X >= other.Min.X) &&
		(w.Min.Y <= other.Max.Y && w.Max.Y >= other.Min.Y) &&
		(w.Min.Z <= other.Max.Z && w.Max.Z >= other.Min.Z)
}

// - Center() Vec3
func (w *WorldBounds) Center() Vec3 {
	return Vec3{
		X: (w.Min.X + w.Max.X) * 0.5,
		Y: (w.Min.Y + w.Max.Y) * 0.5,
		Z: (w.Min.Z + w.Max.Z) * 0.5,
	}
}

// - Size() Vec3
func (w *WorldBounds) Size() Vec3 {
	return Vec3{
		X: w.Max.X - w.Min.X,
		Y: w.Max.Y - w.Min.Y,
		Z: w.Max.Z - w.Min.Z,
	}
}
