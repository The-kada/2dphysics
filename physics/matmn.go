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

func (mat *MatMN) Transpose() *MatMN {
	result := &MatMN{
		elements: make([]VecN, len(mat.elements[0].elements)),
	}

	for index := range result.elements {
		result.elements[index].elements = make([]float64, len(mat.elements))
	}

	for i := 0; i < mat.Rows(); i++ {
		for j := 0; j < mat.Cols(); j++ {
			result.Set(j, i, mat.At(i, j))
		}

	}

	return result
}

func (mat *MatMN) MultiplyMat(matmn *MatMN) (*MatMN, error) {
	if len(mat.elements[0].elements) != len(matmn.elements) {
		return nil, errors.New("mismatched err")
	}

	inputTransposed := matmn.Transpose()

	result := &MatMN{make([]VecN, len(mat.elements))}
	for i := range result.elements {
		result.elements[i] = VecN{elements: make([]float64, len(matmn.elements[0].elements))}
	}

	for i := 0; i < result.Rows(); i++ {
		for j := 0; j < result.Cols(); j++ {
			res, err := mat.elements[i].Dot(&inputTransposed.elements[j])
			if err != nil {
				return nil, err
			}
			result.Set(i, j, res)
		}
	}

	return result, nil
}

// TODO: UNDERSTAND THIS WAS COPIED LINE FOR LINE LOL
func SolveGaussSeidel(matmn *MatMN, vecn *VecN) *VecN {
	n := vecn.Size()

	x := &VecN{elements: make([]float64, n)}

	for iterations := 0; iterations < n; iterations++ {
		for i := 0; i < n; i++ {
			if matmn.At(i, i) != 0 {
				res, err := matmn.elements[i].Dot(x)
				if err != nil {
					panic(err)
				}
				dx := (vecn.At(i) / matmn.At(i, i)) - (res / matmn.At(i, i))
				//if dx == dx {
				x.elements[i] += dx
				//}
			}
		}
	}
	return x
}
