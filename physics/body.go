package physics

type Body struct {
	Postion             Vec2
	velocity            Vec2
	Acceleration        Vec2
	AngularVelocity     float64
	AngularAcceleration float64
	SumTorque           float64
	Mass                float64
	Shape               shape
	Forces              Vec2
	Rotation            float64
	i                   float64
	invI                float64
}

func NewBody(position, velocity, acceleration, forces Vec2, mass, rotation, angularVelocity, angularAcceleration, SumTorque float64, shape shape) *Body {
	var body Body
	body.Postion = position
	body.velocity = velocity
	body.Acceleration = acceleration
	body.AngularAcceleration = angularAcceleration
	body.AngularVelocity = angularVelocity
	body.SumTorque = SumTorque
	body.Mass = mass
	body.Shape = shape
	body.Forces = forces
	body.Rotation = rotation
	body.i = shape.GetMommentOfInertia() * mass
	if body.i != 0 {
		body.invI = 1 / body.i
	} else {
		body.invI = 0
	}

	return &body

}
func (b *Body) CalculateAcceleration() {
	b.Acceleration = b.Forces.Scale(1 / b.Mass)
	b.AngularAcceleration = (b.SumTorque * b.invI)
}

func (b *Body) IntegrateAcceleration(tps float64) {
	time := 0.0
	if tps != 0 {
		time = 1 / tps
	}
	velocity := b.Acceleration.Scale(time)
	b.velocity = b.velocity.Add(&velocity)
	b.AngularVelocity += (b.AngularAcceleration * time)
}

func (b *Body) IntegrateVelocity(tps float64) {
	time := 0.0
	if tps != 0 {
		time = 1 / tps
	}
	posiiton := b.velocity.Scale(1 / tps)
	b.Postion = b.Postion.Add(&posiiton)
	b.Rotation += (b.AngularVelocity * time)
}
