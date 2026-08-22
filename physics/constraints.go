package physics

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
	vec := NewVecN(6)
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
		CachedLambda: NewVecN(1),
	}
}

func (constraint *JointConstraint) GetVelocities() *VecN {
	vec := NewVecN(0, 0, 0, 0, 0, 0)
	vec.elements[0] = constraint.a.velocity[0]
	vec.elements[1] = constraint.a.velocity[1]
	vec.elements[2] = constraint.a.AngularVelocity
	vec.elements[3] = constraint.b.velocity[0]
	vec.elements[4] = constraint.b.velocity[1]
	vec.elements[5] = constraint.b.AngularVelocity
	return vec
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
