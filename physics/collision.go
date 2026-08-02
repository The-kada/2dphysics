package physics

type ContactInfo struct {
	A      *Body
	B      *Body
	Start  Vec2
	End    Vec2
	Depth  float32
	Normal Vec2
}

func IsColliding(a *Body, b *Body) (bool, *ContactInfo) {
	isCircleCircle := a.Shape.GetType() == CircleShape && b.Shape.GetType() == CircleShape
	isPolygonPolygon := a.Shape.GetType() == ConvexPolygonShape || a.Shape.GetType() == RectangleShape && b.Shape.GetType() == ConvexPolygonShape || b.Shape.GetType() == RectangleShape

	if isPolygonPolygon {
		var polygonA, polygonB *ConvexPolygon
		var ok bool

		if polygonA, ok = a.Shape.(*ConvexPolygon); !ok {
			polygonA = &a.Shape.(*Rectangle).ConvexPolygon
		}

		if polygonB, ok = b.Shape.(*ConvexPolygon); !ok {
			polygonB = &b.Shape.(*Rectangle).ConvexPolygon
		}

		sepAB, axisA, pointB := polygonA.FindMinSeperation(polygonB)
		if sepAB >= 0 {
			return false, nil
		}

		sepBA, axisB, pointA := polygonB.FindMinSeperation(polygonA)
		if sepBA >= 0 {
			return false, nil
		}

		contactInfo := &ContactInfo{}
		contactInfo.A = a
		contactInfo.B = b

		if sepAB > sepBA {
			axisANormal := axisA.Normal()
			t := axisANormal.Scale(-float64(sepAB))
			contactInfo.Start = pointB
			contactInfo.End = pointB.Add(&t)
			contactInfo.Depth = -float32(sepAB)
			contactInfo.Normal = axisANormal
		} else {
			axisBNormal := axisB.Normal()
			t := axisBNormal.Scale(-float64(sepBA))
			contactInfo.Start = pointA
			contactInfo.End = pointA.Add(&t)
			contactInfo.Depth = -float32(sepBA)
			contactInfo.Normal = axisBNormal.Scale(-1)
		}
		return true, contactInfo
	}

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

			return true, &ContactInfo{a, b, start, end, float32(depth), vectorAlongCenters.Normal()}
		}
		return false, nil
	} else {
		panic("no!")
	}
}
