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
	Solve()
	PreSolve(dt float64)
}

func (constraint *constraint) GetInvM() *MatMN {
	mat := NewMatMN(6, 6)
	mat.Set(0, 0, constraint.a.invMass)
	mat.Set(1, 1, constraint.a.invMass)
	mat.Set(2, 2, constraint.a.invI)
	mat.Set(3, 3, constraint.b.invMass)
	mat.Set(4, 4, constraint.b.invMass)
	mat.Set(5, 5, constraint.b.invI)
	return mat
}

func (constraint *constraint) GetVelocities() *VecN {
	vec := NewVecN(0, 0, 0, 0, 0, 0)
	vec.elements[0] = constraint.a.velocity[0]
	vec.elements[1] = constraint.a.velocity[1]
	vec.elements[2] = constraint.a.AngularVelocity
	vec.elements[3] = constraint.b.velocity[0]
	vec.elements[4] = constraint.b.velocity[1]
	vec.elements[5] = constraint.b.AngularVelocity
	return vec
}

type JointConstraint struct {
	constraint
	Jacobian     *MatMN
	CachedLambda *VecN
	bias         float64
}

func NewJointConstraint(a, b *Body, anchorPoint Vec2) *JointConstraint {
	return &JointConstraint{
		constraint: constraint{
			a:      a,
			b:      b,
			aPoint: a.WorldSpaceToLocalSpace(&anchorPoint),
			bPoint: b.WorldSpaceToLocalSpace(&anchorPoint),
		},
		Jacobian:     NewMatMN(1, 6),
		CachedLambda: NewVecN(0),
	}
}


func (constraint *JointConstraint) Solve() {
	velocities := constraint.GetVelocities()
	invM := constraint.GetInvM()
	jacobian := constraint.Jacobian
	jt := jacobian.Transpose()

	res, err := jacobian.MultiplyMat(invM)
	if err != nil {
		panic(err)
	}

	lhs, err := res.MultiplyMat(jt)
	if err != nil {
		panic(err)
	}

	rhs, err := jacobian.MultiplyVector(velocities)
	if err != nil {
		panic(err)
	}

	rhs = rhs.Scale(-1)
	rhs.elements[0] -= constraint.bias
	lambda := SolveGaussSeidel(lhs, rhs)
	constraint.CachedLambda, err = constraint.CachedLambda.Add(lambda)
	if err != nil {
		panic(err)
	}

	impulses, err := jt.MultiplyVector(lambda)
	if err != nil {
		panic(err)
	}

	constraint.a.ApplyImpulseLinear(Vec2{impulses.At(0), impulses.At(1)})
	constraint.a.ApplyImpulseAngular(impulses.At(2))

	constraint.b.ApplyImpulseLinear(Vec2{impulses.At(3), impulses.At(4)})
	constraint.b.ApplyImpulseAngular(impulses.At(5))
}

func (constraint *JointConstraint) PreSolve(dt float64) {

	// derive the jacobian
	pa := constraint.a.LocalToWorldSpace(&constraint.aPoint)
	pb := constraint.b.LocalToWorldSpace(&constraint.bPoint)

	ra := pa.Sub(&constraint.a.Postion)
	rb := pb.Sub(&constraint.b.Postion)

	constraint.Jacobian.Zero()

	j1 := (pa.Sub(&pb).Scale(2))
	constraint.Jacobian.Set(0, 0, j1[0])
	constraint.Jacobian.Set(0, 1, j1[1])

	t := pa.Sub(&pb)
	j2 := ra.Cross(&t) * 2
	constraint.Jacobian.Set(0, 2, j2)

	j3 := (pb.Sub(&pa).Scale(2))
	constraint.Jacobian.Set(0, 3, j3[0])
	constraint.Jacobian.Set(0, 4, j3[1])

	t = pb.Sub(&pa)
	j4 := rb.Cross(&t) * 2
	constraint.Jacobian.Set(0, 5, j4)

	jt := constraint.Jacobian.Transpose()
	impulses, err := jt.MultiplyVector(constraint.CachedLambda)
	if err != nil {
		panic(err)
	}

	constraint.a.ApplyImpulseLinear(Vec2{impulses.At(0), impulses.At(1)})
	constraint.a.ApplyImpulseAngular(impulses.At(2))

	constraint.b.ApplyImpulseLinear(Vec2{impulses.At(3), impulses.At(4)})
	constraint.b.ApplyImpulseAngular(impulses.At(5))

	diff := pb.Sub(&pa)
	C := diff.Dot(&diff)
	beta := 0.1
	constraint.bias = (beta / dt) * C
}

type PenetrationConstraint struct {
	constraint
	Jacobian     *MatMN
	CachedLambda *VecN
	bias         float64
	normal       Vec2
	friction     float64
}

func NewPenetrationConstraint(a, b *Body, aCollisionPoint, bCollisionPoint, normal Vec2, friction float64) *PenetrationConstraint {
	return &PenetrationConstraint{
		constraint: constraint{
			a:      a,
			b:      b,
			aPoint: a.WorldSpaceToLocalSpace(&aCollisionPoint),
			bPoint: b.WorldSpaceToLocalSpace(&bCollisionPoint),
		},
		normal:       a.WorldSpaceToLocalSpace(&normal),
		Jacobian:     NewMatMN(2, 6),
		CachedLambda: NewVecN(0, 0),
		friction:     friction,
	}
}

func (constraint *PenetrationConstraint) PreSolve(dt float64) {

	// derive the jacobian
	pa := constraint.a.LocalToWorldSpace(&constraint.aPoint)
	pb := constraint.b.LocalToWorldSpace(&constraint.bPoint)
	n := constraint.a.LocalToWorldSpace(&constraint.normal)

	ra := pa.Sub(&constraint.a.Postion)
	rb := pb.Sub(&constraint.b.Postion)

	constraint.Jacobian.Zero()

	j1 := n.Scale(-1)
	constraint.Jacobian.Set(0, 0, j1[0])
	constraint.Jacobian.Set(0, 1, j1[1])

	j2 := -ra.Cross(&n)
	constraint.Jacobian.Set(0, 2, j2)

	j3 := n
	constraint.Jacobian.Set(0, 3, j3[0])
	constraint.Jacobian.Set(0, 4, j3[1])

	j4 := rb.Cross(&n)
	constraint.Jacobian.Set(0, 5, j4)

	constraint.friction = math.Max(constraint.a.friction, constraint.b.friction)
	if constraint.friction > 0 {
		t := n.Normal().Scale(-1)
		constraint.Jacobian.Set(1, 0, t[0])
		constraint.Jacobian.Set(1, 1, t[1])
		constraint.Jacobian.Set(1, 2, -ra.Cross(&t))

		t = n.Normal()
		constraint.Jacobian.Set(1, 3, t[0])
		constraint.Jacobian.Set(1, 4, t[1])
		constraint.Jacobian.Set(1, 5, rb.Cross(&t))
	}

	jt := constraint.Jacobian.Transpose()

	impulses, err := jt.MultiplyVector(constraint.CachedLambda)
	if err != nil {
		panic(fmt.Errorf("%v matrix dimensions(%d,%d vector length(%d))", err, jt.Rows(), jt.Cols(), constraint.CachedLambda.Size()))
	}

	constraint.a.ApplyImpulseLinear(Vec2{impulses.At(0), impulses.At(1)})
	constraint.a.ApplyImpulseAngular(impulses.At(2))

	constraint.b.ApplyImpulseLinear(Vec2{impulses.At(3), impulses.At(4)})
	constraint.b.ApplyImpulseAngular(impulses.At(5))

	diff := pb.Sub(&pa)
	c := diff.Dot(&j1)
	beta := 0.1
	c = math.Min(0.0, 0.01+c)

	va := constraint.a.velocity.Add(&Vec2{-constraint.a.AngularVelocity * ra[1], -constraint.a.AngularVelocity * ra[0]})
	vb := constraint.b.velocity.Add(&Vec2{-constraint.b.AngularVelocity * rb[1], -constraint.b.AngularVelocity * rb[0]})

	vRelDotNormal := va.Sub(&vb).Dot(&n)
	restitution := math.Min(constraint.a.restitution, constraint.b.restitution)

	constraint.bias = ((beta / dt) * c) + (restitution * vRelDotNormal)
}

func (constraint *PenetrationConstraint) Solve() {
	velocities := constraint.GetVelocities()
	invM := constraint.GetInvM()
	jacobian := constraint.Jacobian
	jt := jacobian.Transpose()

	res, err := jacobian.MultiplyMat(invM)
	if err != nil {
		panic(err)
	}

	lhs, err := res.MultiplyMat(jt)
	if err != nil {
		panic(err)
	}

	rhs, err := jacobian.MultiplyVector(velocities)
	if err != nil {
		panic(err)
	}

	rhs = rhs.Scale(-1)
	rhs.elements[0] -= constraint.bias

	oldLambda := NewVecN(constraint.CachedLambda.elements...)
	lambda := SolveGaussSeidel(lhs, rhs)
	constraint.CachedLambda, err = constraint.CachedLambda.Add(lambda)
	if err != nil {
		panic(err)
	}

	if constraint.CachedLambda.At(0) < 0 {
		constraint.CachedLambda.elements[0] = 0
	}

	if constraint.friction > 0 {
		maxFriction := constraint.CachedLambda.elements[0] * constraint.friction
		clamp := func(x, low, high float64) float64 {
			return math.Max(low, math.Min(x, high))
		}
		constraint.CachedLambda.elements[1] = clamp(constraint.CachedLambda.elements[1], -maxFriction, maxFriction)
	}

	lambda, err = constraint.CachedLambda.Sub(oldLambda)
	if err != nil {
		panic(err)
	}

	impulses, err := jt.MultiplyVector(lambda)
	if err != nil {
		panic(err)
	}

	constraint.a.ApplyImpulseLinear(Vec2{impulses.At(0), impulses.At(1)})
	constraint.a.ApplyImpulseAngular(impulses.At(2))

	constraint.b.ApplyImpulseLinear(Vec2{impulses.At(3), impulses.At(4)})
	constraint.b.ApplyImpulseAngular(impulses.At(5))
}
