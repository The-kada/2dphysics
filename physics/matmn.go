package physics

import "errors"

type matMN struct {
	elements []vecN
}

func newMatMN(row, col int) *matMN {
	mat := &matMN{elements: make([]vecN, row)}
	for index := range mat.elements {
		rowEntry := &vecN{elements: make([]float64, col)}
		mat.elements[index] = *rowEntry
	}
	return mat
}

func (mat *matMN) zero() {
	for _, row := range mat.elements {
		row.zero()
	}
}

func (mat *matMN) at(row, col int) float64 {
	return mat.elements[row].at(col)
}

func (mat *matMN) rows() int {
	return len(mat.elements)
}

func (mat *matMN) cols() int {
	if len(mat.elements) == 0 {
		return 0
	}

	return len(mat.elements[0].elements)
}

func (mat *matMN) set(row, col int, val float64) {
	mat.elements[row].elements[col] = val
}

func (mat *matMN) multiplyVector(vecn *vecN) (*vecN, error) {
	if len(mat.elements[0].elements) != len(vecn.elements) {
		return nil, errors.New("mismatched err")
	}

	result := &vecN{}

	for _, row := range mat.elements {
		val, _ := row.dot(vecn)
		result.elements = append(result.elements, val)
	}

	return result, nil
}

func (mat *matMN) transpose() *matMN {
	result := &matMN{
		elements: make([]vecN, len(mat.elements[0].elements)),
	}

	for index := range result.elements {
		result.elements[index].elements = make([]float64, len(mat.elements))
	}

	for i := 0; i < mat.rows(); i++ {
		for j := 0; j < mat.cols(); j++ {
			result.set(j, i, mat.at(i, j))
		}

	}

	return result
}

func (mat *matMN) multiplyMat(matmn *matMN) (*matMN, error) {
	if len(mat.elements[0].elements) != len(matmn.elements) {
		return nil, errors.New("mismatched err")
	}

	inputTransposed := matmn.transpose()

	result := &matMN{make([]vecN, len(mat.elements))}
	for i := range result.elements {
		result.elements[i] = vecN{elements: make([]float64, len(matmn.elements[0].elements))}
	}

	for i := 0; i < result.rows(); i++ {
		for j := 0; j < result.cols(); j++ {
			res, err := mat.elements[i].dot(&inputTransposed.elements[j])
			if err != nil {
				return nil, err
			}
			result.set(i, j, res)
		}
	}

	return result, nil
}

// TODO: UNDERSTAND THIS WAS COPIED LINE FOR LINE LOL
func solveGaussSeidel(matmn *matMN, vecn *vecN) *vecN {
	n := vecn.size()

	x := &vecN{elements: make([]float64, n)}

	for iterations := 0; iterations < n; iterations++ {
		for i := 0; i < n; i++ {
			if matmn.at(i, i) != 0 {
				res, err := matmn.elements[i].dot(x)
				if err != nil {
					panic(err)
				}
				dx := (vecn.at(i) / matmn.at(i, i)) - (res / matmn.at(i, i))
				//if dx == dx {
				x.elements[i] += dx
				//}
			}
		}
	}
	return x
}
