package physics

import "errors"

type vecN struct {
	elements []float64
}

func newVecN(elements ...float64) *vecN {
	vecN := &vecN{}
	for i := 0; i < len(elements); i++ {
		vecN.elements = append(vecN.elements, elements[i])
	}

	return vecN
}

func (vecn *vecN) add(v2 *vecN) (*vecN, error) {
	if len(v2.elements) != len(vecn.elements) {
		return nil, errors.New("size mismatch")
	}

	newVec := &vecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]+v2.elements[i])
	}

	return newVec, nil
}

func (vecn *vecN) sub(v2 *vecN) (*vecN, error) {
	if len(v2.elements) != len(vecn.elements) {
		return nil, errors.New("size mismatch")
	}

	newVec := &vecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]-v2.elements[i])
	}

	return newVec, nil
}

func (vecn *vecN) scale(n float64) *vecN {

	newVec := &vecN{}

	for i := 0; i < len(vecn.elements); i++ {
		newVec.elements = append(newVec.elements, vecn.elements[i]*n)
	}

	return newVec
}

func (vecn *vecN) dot(v2 *vecN) (float64, error) {
	if len(v2.elements) != len(vecn.elements) {
		return 0, errors.New("size mismatch")
	}

	sum := 0.0

	for i := 0; i < len(vecn.elements); i++ {
		sum += (vecn.elements[i] * v2.elements[i])
	}

	return sum, nil
}

func (vecn *vecN) at(index int) float64 {
	return vecn.elements[index]
}

func (vecn *vecN) size() int {
	return len(vecn.elements)
}

func (vecn *vecN) zero() {
	for i := 0; i < len(vecn.elements); i++ {
		vecn.elements[i] = 0
	}
}
