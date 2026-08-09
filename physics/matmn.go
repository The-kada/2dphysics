package physics

import "errors"

type MatMN struct {
	elements []VecN
}

func NewMatMN(row, col int) *MatMN {
	mat := &MatMN{elements: make([]VecN, row)}
	for index := range mat.elements {
		rowEntry := &VecN{elements: make([]float64, col)}
		mat.elements[index] = *rowEntry
	}
	return mat
}

func (mat *MatMN) Zero() {
	for _, row := range mat.elements {
		row.Zero()
	}
}

func (mat *MatMN) At(row, col int) float64 {
	return mat.elements[row].At(col)
}

func (mat *MatMN) Rows() int {
	return len(mat.elements)
}

func (mat *MatMN) Cols() int {
	if len(mat.elements) == 0 {
		return 0
	}

	return len(mat.elements[0].elements)
}

func (mat *MatMN) Set(row, col int, val float64) {
	mat.elements[row].elements[col] = val
}

func (mat *MatMN) MultiplyVector(vecn *VecN) (*VecN, error) {
	if len(mat.elements[0].elements) != len(vecn.elements) {
		return nil, errors.New("mismatched err")
	}

	result := &VecN{}

	for _, row := range mat.elements {
		val, _ := row.Dot(vecn)
		result.elements = append(result.elements, val)
	}

	return result, nil
}
