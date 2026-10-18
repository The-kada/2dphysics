package physics

import (
	"fmt"
	"math"
)

type constraint struct {
	a      *Body
	b      *Body
	aPoint Vec2
	bPoint Vec2
}

type Constraint interface {
	solve()
	preSolve(dt float64)
}

func (constraint *constraint) getInvM() *matMN {
	mat := newMatMN(6, 6)
	mat.set(0, 0, constraint.a.invMass)
	mat.set(1, 1, constraint.a.invMass)
	mat.set(2, 2, constraint.a.invI)
	mat.set(3, 3, constraint.b.invMass)
	mat.set(4, 4, constraint.b.invMass)
	mat.set(5, 5, constraint.b.invI)
	return mat
}

func (constraint *constraint) getVelocities() *vecN {
	vec := newVecN(0, 0, 0, 0, 0, 0)
	vec.elements[0] = constraint.a.velocity[0]
	vec.elements[1] = constraint.a.velocity[1]
	vec.elements[2] = constraint.a.AngularVelocity
	vec.elements[3] = constraint.b.velocity[0]
	vec.elements[4] = constraint.b.velocity[1]
	vec.elements[5] = constraint.b.AngularVelocity
	return vec
}

type jointConstraint struct {
	constraint
	Jacobian     *matMN
	CachedLambda *vecN
	bias         float64
}

func NewJointConstraint(a, b *Body, anchorPoint Vec2) *jointConstraint {
	return &jointConstraint{
		constraint: constraint{
			a:      a,
			b:      b,
			aPoint: a.worldSpaceToLocalSpace(&anchorPoint),
			bPoint: b.worldSpaceToLocalSpace(&anchorPoint),
		},
		Jacobian:     newMatMN(1, 6),
		CachedLambda: newVecN(0),
	}
}

func (constraint *jointConstraint) solve() {
	velocities := constraint.getVelocities()
	invM := constraint.getInvM()
	jacobian := constraint.Jacobian
	jt := jacobian.transpose()

	res, err := jacobian.multiplyMat(invM)
	if err != nil {
		panic(err)
	}

	lhs, err := res.multiplyMat(jt)
	if err != nil {
		panic(err)
	}

	rhs, err := jacobian.multiplyVector(velocities)
	if err != nil {
		panic(err)
	}

	rhs = rhs.scale(-1)
	rhs.elements[0] -= constraint.bias
	lambda := solveGaussSeidel(lhs, rhs)
	constraint.CachedLambda, err = constraint.CachedLambda.add(lambda)
	if err != nil {
		panic(err)
	}

	impulses, err := jt.multiplyVector(lambda)
	if err != nil {
		panic(err)
	}

	constraint.a.applyImpulseLinear(Vec2{impulses.at(0), impulses.at(1)})
	constraint.a.applyImpulseAngular(impulses.at(2))

	constraint.b.applyImpulseLinear(Vec2{impulses.at(3), impulses.at(4)})
	constraint.b.applyImpulseAngular(impulses.at(5))
}

func (constraint *jointConstraint) preSolve(dt float64) {

	// derive the jacobian
	pa := constraint.a.localToWorldSpace(&constraint.aPoint)
	pb := constraint.b.localToWorldSpace(&constraint.bPoint)

	ra := pa.Sub(&constraint.a.Postion)
	rb := pb.Sub(&constraint.b.Postion)

	constraint.Jacobian.zero()

	j1 := (pa.Sub(&pb).Scale(2))
	constraint.Jacobian.set(0, 0, j1[0])
	constraint.Jacobian.set(0, 1, j1[1])

	t := pa.Sub(&pb)
	j2 := ra.Cross(&t) * 2
	constraint.Jacobian.set(0, 2, j2)

	j3 := (pb.Sub(&pa).Scale(2))
	constraint.Jacobian.set(0, 3, j3[0])
	constraint.Jacobian.set(0, 4, j3[1])

	t = pb.Sub(&pa)
	j4 := rb.Cross(&t) * 2
	constraint.Jacobian.set(0, 5, j4)

	jt := constraint.Jacobian.transpose()
	impulses, err := jt.multiplyVector(constraint.CachedLambda)
	if err != nil {
		panic(err)
	}

	constraint.a.applyImpulseLinear(Vec2{impulses.at(0), impulses.at(1)})
	constraint.a.applyImpulseAngular(impulses.at(2))

	constraint.b.applyImpulseLinear(Vec2{impulses.at(3), impulses.at(4)})
	constraint.b.applyImpulseAngular(impulses.at(5))

	diff := pb.Sub(&pa)
	C := diff.Dot(&diff)
	beta := 0.1
	constraint.bias = (beta / dt) * C
}

type penetrationConstraint struct {
	constraint
	Jacobian     *matMN
	CachedLambda *vecN
	bias         float64
	normal       Vec2
	friction     float64
}

func NewPenetrationConstraint(a, b *Body, aCollisionPoint, bCollisionPoint, normal Vec2, friction float64) *penetrationConstraint {
	return &penetrationConstraint{
		constraint: constraint{
			a:      a,
			b:      b,
			aPoint: a.worldSpaceToLocalSpace(&aCollisionPoint),
			bPoint: b.worldSpaceToLocalSpace(&bCollisionPoint),
		},
		normal:       a.worldSpaceToLocalSpace(&normal),
		Jacobian:     newMatMN(2, 6),
		CachedLambda: newVecN(0, 0),
		friction:     friction,
	}
}

func (constraint *penetrationConstraint) preSolve(dt float64) {

	// derive the jacobian
	pa := constraint.a.localToWorldSpace(&constraint.aPoint)
	pb := constraint.b.localToWorldSpace(&constraint.bPoint)
	n := constraint.a.localToWorldSpace(&constraint.normal)

	ra := pa.Sub(&constraint.a.Postion)
	rb := pb.Sub(&constraint.b.Postion)

	constraint.Jacobian.zero()

	j1 := n.Scale(-1)
	constraint.Jacobian.set(0, 0, j1[0])
	constraint.Jacobian.set(0, 1, j1[1])

	j2 := -ra.Cross(&n)
	constraint.Jacobian.set(0, 2, j2)

	j3 := n
	constraint.Jacobian.set(0, 3, j3[0])
	constraint.Jacobian.set(0, 4, j3[1])

	j4 := rb.Cross(&n)
	constraint.Jacobian.set(0, 5, j4)

	constraint.friction = math.Max(constraint.a.friction, constraint.b.friction)
	if constraint.friction > 0 {
		t := n.Normal().Scale(-1)
		constraint.Jacobian.set(1, 0, t[0])
		constraint.Jacobian.set(1, 1, t[1])
		constraint.Jacobian.set(1, 2, -ra.Cross(&t))

		t = n.Normal()
		constraint.Jacobian.set(1, 3, t[0])
		constraint.Jacobian.set(1, 4, t[1])
		constraint.Jacobian.set(1, 5, rb.Cross(&t))
	}

	jt := constraint.Jacobian.transpose()

	impulses, err := jt.multiplyVector(constraint.CachedLambda)
	if err != nil {
		panic(fmt.Errorf("%v matrix dimensions(%d,%d vector length(%d))", err, jt.rows(), jt.cols(), constraint.CachedLambda.size()))
	}

	constraint.a.applyImpulseLinear(Vec2{impulses.at(0), impulses.at(1)})
	constraint.a.applyImpulseAngular(impulses.at(2))

	constraint.b.applyImpulseLinear(Vec2{impulses.at(3), impulses.at(4)})
	constraint.b.applyImpulseAngular(impulses.at(5))

	diff := pb.Sub(&pa)
	c := diff.Dot(&j1)
	beta := 1.0

	va := constraint.a.velocity.Add(&Vec2{-constraint.a.AngularVelocity * ra[1], constraint.a.AngularVelocity * ra[0]})
	vb := constraint.b.velocity.Add(&Vec2{-constraint.b.AngularVelocity * rb[1], constraint.b.AngularVelocity * rb[0]})

	vRelDotNormal := va.Sub(&vb).Dot(&n)
	restitution := math.Min(constraint.a.restitution, constraint.b.restitution)

	constraint.bias = ((beta / dt) * c) + (restitution * vRelDotNormal)
}

func (constraint *penetrationConstraint) solve() {
	velocities := constraint.getVelocities()
	invM := constraint.getInvM()
	jacobian := constraint.Jacobian
	jt := jacobian.transpose()

	res, err := jacobian.multiplyMat(invM)
	if err != nil {
		panic(err)
	}

	lhs, err := res.multiplyMat(jt)
	if err != nil {
		panic(err)
	}

	rhs, err := jacobian.multiplyVector(velocities)
	if err != nil {
		panic(err)
	}

	rhs = rhs.scale(-1)
	rhs.elements[0] -= constraint.bias

	oldLambda := newVecN(constraint.CachedLambda.elements...)
	lambda := solveGaussSeidel(lhs, rhs)
	constraint.CachedLambda, err = constraint.CachedLambda.add(lambda)
	if err != nil {
		panic(err)
	}

	if constraint.CachedLambda.at(0) < 0 {
		constraint.CachedLambda.elements[0] = 0
	}

	if constraint.friction > 0 {
		maxFriction := constraint.CachedLambda.elements[0] * constraint.friction
		clamp := func(x, low, high float64) float64 {
			return math.Max(low, math.Min(x, high))
		}
		constraint.CachedLambda.elements[1] = clamp(constraint.CachedLambda.elements[1], -maxFriction, maxFriction)
	}

	lambda, err = constraint.CachedLambda.sub(oldLambda)
	if err != nil {
		panic(err)
	}

	impulses, err := jt.multiplyVector(lambda)
	if err != nil {
		panic(err)
	}

	constraint.a.applyImpulseLinear(Vec2{impulses.at(0), impulses.at(1)})
	constraint.a.applyImpulseAngular(impulses.at(2))

	constraint.b.applyImpulseLinear(Vec2{impulses.at(3), impulses.at(4)})
	constraint.b.applyImpulseAngular(impulses.at(5))
}
