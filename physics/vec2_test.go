package physics

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddVec2(t *testing.T) {
	vec1 := Vec2{2, 3}
	vec2 := Vec2{4, 5}
	assert.Equal(t, vec1.Add(&vec2), Vec2{6, 8})
}

func TestMinusVec2(t *testing.T) {
	vec1 := Vec2{2, 3}
	vec2 := Vec2{4, 5}
	assert.Equal(t, vec1.Sub(&vec2), Vec2{-2, -2})
}

func TestScaleVec2(t *testing.T) {
	vec := Vec2{2, 3}
	assert.Equal(t, vec.Scale(10), Vec2{20, 30})
}

func TestRotateVec2(t *testing.T) {
	vec := Vec2{1, 1}
	rotatedVector := vec.Rotate(90)
	assert.True(t, almostEqual(rotatedVector[0], -1, 1e-9))
	assert.True(t, almostEqual(rotatedVector[1], 1, 1e-9))
}

func TestMagnitudeVec2(t *testing.T) {
	vec := Vec2{6, 8}
	assert.Equal(t, vec.Magnitude(), float64(10))
}

func TestMagnitudeSquaredVec2(t *testing.T) {
	vec := Vec2{6, 8}
	assert.Equal(t, vec.MagnitudeSquared(), float64(100))
}

func TestUnitVectorVec2(t *testing.T) {
	vec := Vec2{1, 1}
	magnitude := vec.UnitVector().Magnitude()
	assert.True(t, almostEqual(magnitude, 1, 1e-9))

	zeroVector := Vec2{0, 0}
	assert.Equal(t, zeroVector.UnitVector(), Vec2{0, 0})
}

func TestNormalizeVectorVec2(t *testing.T) {
	vec := Vec2{5, 0}
	vec.Normalize()
	assert.Equal(t, vec, Vec2{1, 0})

}

func TestDotVectorVec2(t *testing.T) {
	vec1 := Vec2{2, 2}
	vec2 := Vec2{4, 5}
	assert.Equal(t, vec1.Dot(&vec2), float64(18))

}

func TestCrossVectorVec2(t *testing.T) {
	vec1 := Vec2{20, 20}
	vec2 := Vec2{4, 5}
	assert.Equal(t, vec1.Cross(&vec2), float64(20))

}

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}
