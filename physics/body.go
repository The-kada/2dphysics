package physics

type Body struct {
	Postion      Vec2
	velocity     Vec2
	Acceleration Vec2
	Mass         float64
	Shape        shape
	Forces       Vec2
}

func (b *Body) CalculateAcceleration() {
	b.Acceleration = b.Forces.Scale(1 / b.Mass)
}

func (b *Body) IntegrateAcceleration(tps float64) {
	velocity := b.Acceleration.Scale(1 / tps)
	b.velocity = b.velocity.Add(&velocity)
}

func (b *Body) IntegratePosition(tps float64) {
	posiiton := b.velocity.Scale(1 / tps)
	b.Postion = b.Postion.Add(&posiiton)
}
