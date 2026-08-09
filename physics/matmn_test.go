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
