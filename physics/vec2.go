package physics

import "math"

type Vec2 [2]float64

func (v1 Vec2) Add(v2 *Vec2) Vec2 {
	return Vec2{v1[0] + v2[0], v1[1] + v2[1]}
}

func (v1 Vec2) Sub(v2 *Vec2) Vec2 {
	return Vec2{v1[0] - v2[0], v1[1] - v2[1]}
}

func (v1 Vec2) Scale(n float64) Vec2 {
	return Vec2{v1[0] * n, v1[1] * n}
}

func (v1 Vec2) Normal() Vec2 {
	v := Vec2{v1[1], -v1[0]}
	v.Normalize()
	return v
}

// Rotate the vector in degrees
func (v1 Vec2) Rotate(angle float64) Vec2 {
	angle = angle * (math.Pi / 180)
	xRotation := v1[0]*math.Cos(angle) - v1[1]*math.Sin(angle)
	yRotation := v1[0]*math.Sin(angle) + v1[1]*math.Cos(angle)
	return Vec2{xRotation, yRotation}
}

func (v1 Vec2) Magnitude() float64 {
	return math.Sqrt(v1[0]*v1[0] + v1[1]*v1[1])
}

func (v1 *Vec2) MagnitudeSquared() float64 {
	return (v1[0]*v1[0] + v1[1]*v1[1])
}

func (v1 Vec2) UnitVector() Vec2 {
	magnitude := v1.Magnitude()
	if magnitude == 0 {
		return Vec2{v1[0], v1[1]}
	}
	return Vec2{v1[0] / magnitude, v1[1] / magnitude}
}

func (v1 *Vec2) Normalize() {
	magnitude := v1.Magnitude()
	if magnitude == 0 {
		return
	}
	v1[0] = v1[0] / magnitude
	v1[1] = v1[1] / magnitude
}

func (v1 Vec2) Dot(v2 *Vec2) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func (v1 Vec2) Cross(v2 *Vec2) float64 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}
