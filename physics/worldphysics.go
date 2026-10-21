package physics

type WorldPhysics struct {
	Bodies      []*Body
	Contacts    []*ContactInfo
	Constraints []Constraint
	BodyForces  map[*Body][]Vec2
}

func (w *WorldPhysics) AddBody(body ...*Body) {
	w.Bodies = append(w.Bodies, body...)
}

func (w *WorldPhysics) AddContactInfo(contactInfo ...*ContactInfo) {
	w.Contacts = append(w.Contacts, contactInfo...)
}

func (w *WorldPhysics) AddConstraint(constraint ...Constraint) {
	w.Constraints = append(w.Constraints, constraint...)
}

func (w *WorldPhysics) UpdatePhysics(deltaTime float64) {
	w.Contacts = make([]*ContactInfo, 0, 0)

	for i := 0; i < len(w.Bodies); i++ {
		for j := i + 1; j < len(w.Bodies); j++ {
			if colliding, contacts := isColliding(w.Bodies[i], w.Bodies[j]); colliding {
				w.Contacts = append(w.Contacts, contacts...)
			}
		}
	}

	w.Constraints = make([]Constraint, 0, 0)

	for _, contact := range w.Contacts {
		constraint := NewPenetrationConstraint(contact.A, contact.B, contact.Start, contact.End, contact.Normal, 0)
		w.Constraints = append(w.Constraints, constraint)
	}

	for _, body := range w.Bodies {
		if forces, present := w.BodyForces[body]; present {
			for _, force := range forces {
				f := force
				body.Forces = body.Forces.Add(&f)
			}
		}

	}

	for _, body := range w.Bodies {
		body.integrateAcceleration(deltaTime)
	}

	for _, constraint := range w.Constraints {
		constraint.preSolve(1.0 / deltaTime)
	}

	for i := 0; i < 5; i++ {
		for _, constraint := range w.Constraints {
			constraint.solve()
		}
	}

	for _, body := range w.Bodies {
		body.integrateVelocity(deltaTime)
	}

}
