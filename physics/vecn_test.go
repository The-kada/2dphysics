package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// VecN is an arbitrary-length float64 vector.
// Tests below are the spec — implement VecN to make them pass.

func TestVecNAdd(t *testing.T) {
	v1 := NewVecN(1, 2, 3)
	v2 := NewVecN(4, 5, 6)
	result, err := v1.Add(v2)
	require.NoError(t, err)
	assert.Equal(t, NewVecN(5, 7, 9), result)
}

func TestVecNAddMismatchedSizes(t *testing.T) {
	v1 := NewVecN(1, 2, 3)
	v2 := NewVecN(4, 5)
	_, err := v1.Add(v2)
	assert.Error(t, err)
}

func TestVecNSub(t *testing.T) {
	v1 := NewVecN(5, 7, 9)
	v2 := NewVecN(1, 2, 3)
	result, err := v1.Sub(v2)
	require.NoError(t, err)
	assert.Equal(t, NewVecN(4, 5, 6), result)
}

func TestVecNSubMismatchedSizes(t *testing.T) {
	v1 := NewVecN(1, 2)
	v2 := NewVecN(1, 2, 3)
	_, err := v1.Sub(v2)
	assert.Error(t, err)
}

func TestVecNScale(t *testing.T) {
	v := NewVecN(2, 3, 4)
	assert.Equal(t, NewVecN(6, 9, 12), v.Scale(3))
}

func TestVecNScaleNegative(t *testing.T) {
	v := NewVecN(2, -3, 4)
	assert.Equal(t, NewVecN(-4, 6, -8), v.Scale(-2))
}

func TestVecNDot(t *testing.T) {
	v1 := NewVecN(1, 2, 3)
	v2 := NewVecN(4, 5, 6)
	result, err := v1.Dot(v2)
	require.NoError(t, err)
	assert.Equal(t, float64(32), result) // 1*4 + 2*5 + 3*6
}

func TestVecNDotMismatchedSizes(t *testing.T) {
	v1 := NewVecN(1, 2)
	v2 := NewVecN(1, 2, 3)
	_, err := v1.Dot(v2)
	assert.Error(t, err)
}

func TestVecNDotOrthogonalVectorsIsZero(t *testing.T) {
	v1 := NewVecN(1, 0, 0)
	v2 := NewVecN(0, 1, 0)
	result, err := v1.Dot(v2)
	require.NoError(t, err)
	assert.Equal(t, float64(0), result)
}

func TestVecNSize(t *testing.T) {
	v := NewVecN(1, 2, 3, 4, 5)
	assert.Equal(t, 5, v.Size())
}

func TestVecNAt(t *testing.T) {
	v := NewVecN(10, 20, 30)
	assert.Equal(t, float64(10), v.At(0))
	assert.Equal(t, float64(20), v.At(1))
	assert.Equal(t, float64(30), v.At(2))
}

func TestVecNZeroSetsAllElementsToZero(t *testing.T) {
	v := NewVecN(1, 2, 3)
	v.Zero()
	assert.Equal(t, NewVecN(0, 0, 0), v)
}

func TestVecNZeroPreservesSize(t *testing.T) {
	v := NewVecN(1, 2, 3, 4, 5)
	v.Zero()
	assert.Equal(t, 5, v.Size())
}

func TestVecNZeroOnAlreadyZeroVector(t *testing.T) {
	v := NewVecN(0, 0, 0)
	v.Zero()
	assert.Equal(t, NewVecN(0, 0, 0), v)
}
