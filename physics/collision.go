package physics

type ContactInfo struct {
	A     *Body
	B     *Body
	Start Vec2
	End   Vec2
	Depth float32
}

func IsColliding(a *Body, b *Body) (bool, *ContactInfo) {
	isCircleCircle := a.Shape.GetType() == CircleShape && b.Shape.GetType() == CircleShape

	if isCircleCircle {
		vectorAlongCenters := b.Postion.Sub(&a.Postion)
		distanceBetweenCircles := vectorAlongCenters.Magnitude()
		aRadius := a.Shape.(*Circle).Radius
		bRadius := b.Shape.(*Circle).Radius
		if distanceBetweenCircles <= aRadius+bRadius {
			vectorAlongCenters.Normalize()
			normalAlongCenters := vectorAlongCenters

			bToA := normalAlongCenters.Scale(bRadius)
			aToB := normalAlongCenters.Scale(aRadius)

			start := b.Postion.Sub(&bToA)
			end := a.Postion.Add(&aToB)
			depth := end.Sub(&start).Magnitude()

			return true, &ContactInfo{a, b, start, end, float32(depth)}
		}
		return false, nil
	} else {
		panic("no!")
	}
}
