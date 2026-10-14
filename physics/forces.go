package physics

const PixelPerMeter = 20
const G = 9.8

func GenerateSpringForce(body *Body, anchor Vec2, k float64, restLength float64) *Vec2 {
	d := body.Postion.Sub(&anchor)
	displacement := d.Magnitude() - float64(restLength)
	direction := d.UnitVector()
	springMagnitude := -k * displacement
	force := direction.Scale(springMagnitude)
	return &force
}
