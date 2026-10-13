package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// VecN is an arbitrary-length float64 vector.
// Tests below are the spec — implement VecN to make them pass.

func TestVecNAdd(t *testing.T) {
	v1 := newVecN(1, 2, 3)
	v2 := newVecN(4, 5, 6)
	result, err := v1.add(v2)
	require.NoError(t, err)
	assert.Equal(t, newVecN(5, 7, 9), result)
}

func TestVecNAddMismatchedSizes(t *testing.T) {
	v1 := newVecN(1, 2, 3)
	v2 := newVecN(4, 5)
	_, err := v1.add(v2)
	assert.Error(t, err)
}

func TestVecNSub(t *testing.T) {
	v1 := newVecN(5, 7, 9)
	v2 := newVecN(1, 2, 3)
	result, err := v1.sub(v2)
	require.NoError(t, err)
	assert.Equal(t, newVecN(4, 5, 6), result)
}

func TestVecNSubMismatchedSizes(t *testing.T) {
	v1 := newVecN(1, 2)
	v2 := newVecN(1, 2, 3)
	_, err := v1.sub(v2)
	assert.Error(t, err)
}

func TestVecNScale(t *testing.T) {
	v := newVecN(2, 3, 4)
	assert.Equal(t, newVecN(6, 9, 12), v.scale(3))
}

func TestVecNScaleNegative(t *testing.T) {
	v := newVecN(2, -3, 4)
	assert.Equal(t, newVecN(-4, 6, -8), v.scale(-2))
}

func TestVecNDot(t *testing.T) {
	v1 := newVecN(1, 2, 3)
	v2 := newVecN(4, 5, 6)
	result, err := v1.dot(v2)
	require.NoError(t, err)
	assert.Equal(t, float64(32), result) // 1*4 + 2*5 + 3*6
}

func TestVecNDotMismatchedSizes(t *testing.T) {
	v1 := newVecN(1, 2)
	v2 := newVecN(1, 2, 3)
	_, err := v1.dot(v2)
	assert.Error(t, err)
}

func TestVecNDotOrthogonalVectorsIsZero(t *testing.T) {
	v1 := newVecN(1, 0, 0)
	v2 := newVecN(0, 1, 0)
	result, err := v1.dot(v2)
	require.NoError(t, err)
	assert.Equal(t, float64(0), result)
}

func TestVecNSize(t *testing.T) {
	v := newVecN(1, 2, 3, 4, 5)
	assert.Equal(t, 5, v.size())
}

func TestVecNAt(t *testing.T) {
	v := newVecN(10, 20, 30)
	assert.Equal(t, float64(10), v.at(0))
	assert.Equal(t, float64(20), v.at(1))
	assert.Equal(t, float64(30), v.at(2))
}

func TestVecNZeroSetsAllElementsToZero(t *testing.T) {
	v := newVecN(1, 2, 3)
	v.zero()
	assert.Equal(t, newVecN(0, 0, 0), v)
}

func TestVecNZeroPreservesSize(t *testing.T) {
	v := newVecN(1, 2, 3, 4, 5)
	v.zero()
	assert.Equal(t, 5, v.size())
}

func TestVecNZeroOnAlreadyZeroVector(t *testing.T) {
	v := newVecN(0, 0, 0)
	v.zero()
	assert.Equal(t, newVecN(0, 0, 0), v)
}
