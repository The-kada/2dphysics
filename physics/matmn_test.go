package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatMNDimensions(t *testing.T) {
	mat := NewMatMN(3, 4)
	assert.Equal(t, 3, mat.Rows())
	assert.Equal(t, 4, mat.Cols())
}

func TestNewMatMNInitialisedToZero(t *testing.T) {
	mat := NewMatMN(2, 3)
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			assert.Equal(t, float64(0), mat.At(r, c))
		}
	}
}

func TestMatMNAt(t *testing.T) {
	mat := NewMatMN(2, 2)
	mat.Set(0, 0, 1)
	mat.Set(0, 1, 2)
	mat.Set(1, 0, 3)
	mat.Set(1, 1, 4)
	assert.Equal(t, float64(1), mat.At(0, 0))
	assert.Equal(t, float64(2), mat.At(0, 1))
	assert.Equal(t, float64(3), mat.At(1, 0))
	assert.Equal(t, float64(4), mat.At(1, 1))
}

func TestMatMNZeroSetsAllToZero(t *testing.T) {
	mat := NewMatMN(2, 3)
	mat.Set(0, 0, 5)
	mat.Set(1, 2, 9)
	mat.Zero()
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			assert.Equal(t, float64(0), mat.At(r, c))
		}
	}
}

func TestMatMNZeroPreservesDimensions(t *testing.T) {
	mat := NewMatMN(4, 5)
	mat.Zero()
	assert.Equal(t, 4, mat.Rows())
	assert.Equal(t, 5, mat.Cols())
}

func TestMatMNMultiplyVector(t *testing.T) {
	// [1 2 3] * [1]   [14]
	// [4 5 6]   [2] = [32]
	//           [3]
	mat := NewMatMN(2, 3)
	mat.Set(0, 0, 1); mat.Set(0, 1, 2); mat.Set(0, 2, 3)
	mat.Set(1, 0, 4); mat.Set(1, 1, 5); mat.Set(1, 2, 6)
	v := NewVecN(1, 2, 3)
	result, err := mat.MultiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, float64(14), result.At(0))
	assert.Equal(t, float64(32), result.At(1))
}

func TestMatMNMultiplyVectorResultSize(t *testing.T) {
	mat := NewMatMN(3, 2)
	v := NewVecN(1, 2)
	result, err := mat.MultiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, 3, result.Size())
}

func TestMatMNMultiplyVectorMismatchedSizes(t *testing.T) {
	mat := NewMatMN(2, 3)
	v := NewVecN(1, 2) // cols=3, vec size=2
	_, err := mat.MultiplyVector(v)
	assert.Error(t, err)
}

func TestMatMNMultiplyVectorIdentity(t *testing.T) {
	mat := NewMatMN(3, 3)
	mat.Set(0, 0, 1); mat.Set(1, 1, 1); mat.Set(2, 2, 1)
	v := NewVecN(4, 5, 6)
	result, err := mat.MultiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, float64(4), result.At(0))
	assert.Equal(t, float64(5), result.At(1))
	assert.Equal(t, float64(6), result.At(2))
}
