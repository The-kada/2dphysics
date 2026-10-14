package physics

import (
	"math"
)

type Body struct {
	Postion             Vec2
	velocity            Vec2
	Acceleration        Vec2
	AngularVelocity     float64
	AngularAcceleration float64
	SumTorque           float64
	Mass                float64
	invMass             float64
	Shape               shape
	Forces              Vec2
	Rotation            float64
	i                   float64
	invI                float64
	friction            float64
	restitution         float64
}

func NewBody(position, velocity, acceleration, forces Vec2, mass, rotation, angularVelocity, angularAcceleration, SumTorque, friciton, restitution float64, shape shape) *Body {
	var body Body
	body.Postion = position
	body.velocity = velocity
	body.Acceleration = acceleration
	body.AngularAcceleration = angularAcceleration
	body.AngularVelocity = angularVelocity
	body.SumTorque = SumTorque
	body.Mass = mass
	body.friction = friciton
	body.restitution = restitution
	if body.Mass != 0 {
		body.invMass = 1 / body.Mass
	} else {
		body.invMass = 0
	}
	body.Shape = shape
	body.Forces = forces
	body.Rotation = rotation
	body.i = shape.getMommentOfInertia() * mass
	if body.i != 0 {
		body.invI = 1 / body.i
	} else {
		body.invI = 0
	}

	return &body

}

func (b *Body) isStatic() bool {
	return b.invMass == 0
}

func (b *Body) integrateAcceleration(tps float64) {
	if b.isStatic() {
		return
	}

	time := 0.0
	if tps != 0 {
		time = 1 / tps
	}

	b.Acceleration = b.Forces.Scale(1 / b.Mass)
	b.AngularAcceleration = (b.SumTorque * b.invI)

	velocity := b.Acceleration.Scale(time)
	b.velocity = b.velocity.Add(&velocity)
	b.AngularVelocity += (b.AngularAcceleration * time)
	b.SumTorque = 0
	b.Forces[0], b.Forces[1] = 0, 0
}

func (b *Body) integrateVelocity(tps float64) {
	if b.isStatic() {
		return
	}
	time := 0.0
	if tps != 0 {
		time = 1 / tps
	}
	posiiton := b.velocity.Scale(time)
	b.Postion = b.Postion.Add(&posiiton)
	b.Rotation += (b.AngularVelocity * time)
}

func (b *Body) worldSpaceToLocalSpace(point *Vec2) Vec2 {
	translatedX := point[0] - b.Postion[0]
	translatedY := point[1] - b.Postion[1]

	rotatedX := math.Cos(-b.Rotation)*translatedX - math.Sin(-b.Rotation)*translatedY
	rotatedY := math.Cos(-b.Rotation)*translatedY + math.Sin(-b.Rotation)*translatedX

	return Vec2{rotatedX, rotatedY}
}

func (b *Body) localToWorldSpace(point *Vec2) Vec2 {
	rotated := point.Rotate(b.Rotation)
	return rotated.Add(&b.Postion)
}

func (b *Body) applyImpulseLinear(impulse Vec2) {
	if b.isStatic() {
		return
	}

	x := impulse.Scale(b.invMass)
	b.velocity = b.velocity.Add(&x)
}

func (b *Body) applyImpulseAngular(impulse float64) {
	if b.isStatic() {
		return
	}
	b.AngularVelocity += impulse * b.invI
}

func (b *Body) applyImpulseAtPoint(impulse Vec2, radius Vec2) {
	if b.isStatic() {
		return
	}
	x := impulse.Scale(b.invMass)
	b.velocity = b.velocity.Add(&x)

	b.AngularVelocity += radius.Cross(&impulse) * b.invI
}
