package physics

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
