package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMatMNDimensions(t *testing.T) {
	mat := newMatMN(3, 4)
	assert.Equal(t, 3, mat.rows())
	assert.Equal(t, 4, mat.cols())
}

func TestNewMatMNInitialisedToZero(t *testing.T) {
	mat := newMatMN(2, 3)
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			assert.Equal(t, float64(0), mat.at(r, c))
		}
	}
}

func TestMatMNAt(t *testing.T) {
	mat := newMatMN(2, 2)
	mat.set(0, 0, 1)
	mat.set(0, 1, 2)
	mat.set(1, 0, 3)
	mat.set(1, 1, 4)
	assert.Equal(t, float64(1), mat.at(0, 0))
	assert.Equal(t, float64(2), mat.at(0, 1))
	assert.Equal(t, float64(3), mat.at(1, 0))
	assert.Equal(t, float64(4), mat.at(1, 1))
}

func TestMatMNZeroSetsAllToZero(t *testing.T) {
	mat := newMatMN(2, 3)
	mat.set(0, 0, 5)
	mat.set(1, 2, 9)
	mat.zero()
	for r := 0; r < 2; r++ {
		for c := 0; c < 3; c++ {
			assert.Equal(t, float64(0), mat.at(r, c))
		}
	}
}

func TestMatMNZeroPreservesDimensions(t *testing.T) {
	mat := newMatMN(4, 5)
	mat.zero()
	assert.Equal(t, 4, mat.rows())
	assert.Equal(t, 5, mat.cols())
}

func TestMatMNMultiplyVector(t *testing.T) {
	// [1 2 3] * [1]   [14]
	// [4 5 6]   [2] = [32]
	//           [3]
	mat := newMatMN(2, 3)
	mat.set(0, 0, 1)
	mat.set(0, 1, 2)
	mat.set(0, 2, 3)
	mat.set(1, 0, 4)
	mat.set(1, 1, 5)
	mat.set(1, 2, 6)
	v := newVecN(1, 2, 3)
	result, err := mat.multiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, float64(14), result.at(0))
	assert.Equal(t, float64(32), result.at(1))
}

func TestMatMNMultiplyVectorResultSize(t *testing.T) {
	mat := newMatMN(3, 2)
	v := newVecN(1, 2)
	result, err := mat.multiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, 3, result.size())
}

func TestMatMNMultiplyVectorMismatchedSizes(t *testing.T) {
	mat := newMatMN(2, 3)
	v := newVecN(1, 2) // cols=3, vec size=2
	_, err := mat.multiplyVector(v)
	assert.Error(t, err)
}

func TestMatMNTransposeDimensions(t *testing.T) {
	mat := newMatMN(2, 3)
	transposed := mat.transpose()
	assert.Equal(t, 3, transposed.rows())
	assert.Equal(t, 2, transposed.cols())
}

func TestMatMNTransposeValues(t *testing.T) {
	// [1 2 3]T   [1 4]
	// [4 5 6]  = [2 5]
	//            [3 6]
	mat := newMatMN(2, 3)
	mat.set(0, 0, 1)
	mat.set(0, 1, 2)
	mat.set(0, 2, 3)
	mat.set(1, 0, 4)
	mat.set(1, 1, 5)
	mat.set(1, 2, 6)
	t2 := mat.transpose()
	assert.Equal(t, float64(1), t2.at(0, 0))
	assert.Equal(t, float64(4), t2.at(0, 1))
	assert.Equal(t, float64(2), t2.at(1, 0))
	assert.Equal(t, float64(5), t2.at(1, 1))
	assert.Equal(t, float64(3), t2.at(2, 0))
	assert.Equal(t, float64(6), t2.at(2, 1))
}

func TestMatMNTransposeSquare(t *testing.T) {
	mat := newMatMN(2, 2)
	mat.set(0, 0, 1)
	mat.set(0, 1, 2)
	mat.set(1, 0, 3)
	mat.set(1, 1, 4)
	t2 := mat.transpose()
	assert.Equal(t, float64(1), t2.at(0, 0))
	assert.Equal(t, float64(3), t2.at(0, 1))
	assert.Equal(t, float64(2), t2.at(1, 0))
	assert.Equal(t, float64(4), t2.at(1, 1))
}

func TestMatMNTransposeDoubleIsOriginal(t *testing.T) {
	mat := newMatMN(2, 3)
	mat.set(0, 0, 1)
	mat.set(0, 1, 2)
	mat.set(0, 2, 3)
	mat.set(1, 0, 4)
	mat.set(1, 1, 5)
	mat.set(1, 2, 6)
	tt := mat.transpose().transpose()
	for r := range 2 {
		for c := range 3 {
			assert.Equal(t, mat.at(r, c), tt.at(r, c))
		}
	}
}

func TestMatMNMultiplyVectorIdentity(t *testing.T) {
	mat := newMatMN(3, 3)
	mat.set(0, 0, 1)
	mat.set(1, 1, 1)
	mat.set(2, 2, 1)
	v := newVecN(4, 5, 6)
	result, err := mat.multiplyVector(v)
	assert.NoError(t, err)
	assert.Equal(t, float64(4), result.at(0))
	assert.Equal(t, float64(5), result.at(1))
	assert.Equal(t, float64(6), result.at(2))
}

func TestMatMNMultiplyMat(t *testing.T) {
	// [1 2] * [5 6]   [19 22]
	// [3 4]   [7 8] = [43 50]
	a := newMatMN(2, 2)
	a.set(0, 0, 1)
	a.set(0, 1, 2)
	a.set(1, 0, 3)
	a.set(1, 1, 4)

	b := newMatMN(2, 2)
	b.set(0, 0, 5)
	b.set(0, 1, 6)
	b.set(1, 0, 7)
	b.set(1, 1, 8)

	result, err := a.multiplyMat(b)
	assert.NoError(t, err)
	assert.Equal(t, float64(19), result.at(0, 0))
	assert.Equal(t, float64(22), result.at(0, 1))
	assert.Equal(t, float64(43), result.at(1, 0))
	assert.Equal(t, float64(50), result.at(1, 1))
}

func TestMatMNMultiplyMatResultDimensions(t *testing.T) {
	// (2x3) * (3x4) = (2x4)
	a := newMatMN(2, 3)
	b := newMatMN(3, 4)
	result, err := a.multiplyMat(b)
	assert.NoError(t, err)
	assert.Equal(t, 2, result.rows())
	assert.Equal(t, 4, result.cols())
}

func TestMatMNMultiplyMatMismatchedSizes(t *testing.T) {
	a := newMatMN(2, 3)
	b := newMatMN(4, 2) // a.Cols=3 != b.Rows=4
	_, err := a.multiplyMat(b)
	assert.Error(t, err)
}

func TestMatMNMultiplyMatIdentity(t *testing.T) {
	a := newMatMN(2, 2)
	a.set(0, 0, 3)
	a.set(0, 1, 7)
	a.set(1, 0, 2)
	a.set(1, 1, 5)

	identity := newMatMN(2, 2)
	identity.set(0, 0, 1)
	identity.set(1, 1, 1)

	result, err := a.multiplyMat(identity)
	assert.NoError(t, err)
	for r := range 2 {
		for c := range 2 {
			assert.Equal(t, a.at(r, c), result.at(r, c))
		}
	}
}
