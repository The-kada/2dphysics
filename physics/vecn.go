package physics

import "errors"

type VecN struct {
	elements []float64
}

func NewVecN(elements ...float64) *VecN {
	vecN := &VecN{}
	for i := 0; i < len(elements); i++ {
		vecN.elements = append(vecN.elements, elements[i])
	}

	return vecN
}

func (vecn *VecN) Add(v2 *VecN) (*VecN, error) {
	if len(v2.elements) != len(vecn.elements) {
		return nil, errors.New("size mismatch")
	}

	newVec := &VecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]+v2.elements[i])
	}

	return newVec, nil
}

func (vecn *VecN) Sub(v2 *VecN) (*VecN, error) {
	if len(v2.elements) != len(vecn.elements) {
		return nil, errors.New("size mismatch")
	}

	newVec := &VecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]-v2.elements[i])
	}

	return newVec, nil
}

func (vecn *VecN) Scale(n float64) *VecN {

	newVec := &VecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]*n)
	}

	return newVec
}

func (vecn *VecN) Dot(v2 *VecN) (float64, error) {
	if len(v2.elements) != len(vecn.elements) {
		return 0, errors.New("size mismatch")
	}

	sum := 0.0

	for i := 0; i < len(vecn.elements); i++ {
		sum += (vecn.elements[i] * v2.elements[i])
	}

	return sum, nil
}

func (vecn *VecN) At(index int) float64 {
	return vecn.elements[index]
}

func (vecn *VecN) Size() int {
	return len(vecn.elements)
}

func (vecn *VecN) Zero() {
	for i := 0; i < len(vecn.elements); i++ {
		vecn.elements[i] = 0
	}
}
