package physics

import "math"

type ContactInfo struct {
	A      *body
	B      *body
	Start  Vec2
	End    Vec2
	Depth  float32
	Normal Vec2
}

func isColliding(a *body, b *body) (bool, *ContactInfo) {
	isCircleCircle := a.Shape.GetType() == CircleShape && b.Shape.GetType() == CircleShape
	isPolygonPolygon := (a.Shape.GetType() == ConvexPolygonShape || a.Shape.GetType() == RectangleShape) && (b.Shape.GetType() == ConvexPolygonShape || b.Shape.GetType() == RectangleShape)
	isPolygonCircle := (a.Shape.GetType() == CircleShape && b.Shape.GetType() == RectangleShape || b.Shape.GetType() == ConvexPolygonShape) || (b.Shape.GetType() == CircleShape && a.Shape.GetType() == RectangleShape || a.Shape.GetType() == ConvexPolygonShape)

	if isPolygonPolygon {
		return checkPolygonPolygonCollision(a, b)
	}

	if isPolygonCircle {
		return checkCirclePolygonCollision(a, b)

	}

	if isCircleCircle {
		return checkCircleCircleCollision(a, b)
	}

	panic("impossible!")

}

func checkCircleCircleCollision(a *body, b *body) (bool, *ContactInfo) {
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

		return true, &ContactInfo{a, b, start, end, float32(depth), normalAlongCenters}
	}

	return false, nil
}

func checkCirclePolygonCollision(a *body, b *body) (bool, *ContactInfo) {

	var contact ContactInfo
	var polygon *ConvexPolygon
	var circle *Circle

	switch a.Shape.GetType() {
	case ConvexPolygonShape:
		contact.A = a
		polygon = a.Shape.(*ConvexPolygon)
	case RectangleShape:
		contact.A = a
		polygon = &a.Shape.(*Rectangle).ConvexPolygon
	case CircleShape:
		contact.B = a
		circle = a.Shape.(*Circle)
	}

	switch b.Shape.GetType() {
	case ConvexPolygonShape:
		contact.A = b
		polygon = b.Shape.(*ConvexPolygon)
	case RectangleShape:
		contact.A = b
		polygon = &b.Shape.(*Rectangle).ConvexPolygon
	case CircleShape:
		contact.B = b
		circle = b.Shape.(*Circle)
	}

	distanceCircleEdge := -math.MaxFloat64
	isOutside := false
	var currentVertex Vec2
	var nextVertex Vec2
	for i := 0; i < len(polygon.WorldVertices); i++ {
		currentVertex = polygon.WorldVertices[i]
		nextVertex = polygon.WorldVertices[(i+1)%len(polygon.WorldVertices)]
		normal := polygon.edgeAtIndex(i).Normal().UnitVector()

		currentVertexToCircleCenter := contact.B.Postion.Sub(&currentVertex)

		proj := normal.Dot(&currentVertexToCircleCenter)
		if proj <= 0 {
			if proj > distanceCircleEdge {
				distanceCircleEdge = proj
				contact.Normal = normal
				contact.Depth = float32(circle.Radius) + (-float32(distanceCircleEdge))
				rad := contact.Normal.Scale(-circle.Radius)
				contact.Start = contact.B.Postion.Add(&rad)
				normalScaledBydepth := contact.Normal.Scale(float64(contact.Depth))
				contact.End = contact.Start.Add(&normalScaledBydepth)
			}
		} else {
			isOutside = true
			distanceCircleEdge = proj
			break
		}
	}

	if !isOutside {
		return true, &contact
	}

	currentVertexVectorCircle := contact.B.Postion.Sub(&currentVertex)
	nextVertexVectorCircle := contact.B.Postion.Sub(&nextVertex)
	currVertexToNextVertexVector := nextVertex.Sub(&currentVertex)
	nextVertexToCurrVertexVector := currVertexToNextVertexVector.Scale(-1)

	if currentVertexVectorCircle.Dot(&currVertexToNextVertexVector) < 0 {
		//region A
		if currentVertexVectorCircle.Magnitude() > circle.Radius {
			return false, nil
		}

		contact.Normal = currentVertexVectorCircle.UnitVector()
		normalScaledRadius := contact.Normal.Scale(circle.Radius)
		contact.Start = contact.B.Postion.Sub(&normalScaledRadius)
		contact.Depth = float32(circle.Radius) - float32(currentVertexVectorCircle.Magnitude())
		depthScaledNormal := contact.Normal.Scale(float64(contact.Depth))
		contact.End = contact.Start.Add(&depthScaledNormal)
		return true, &contact

	} else if nextVertexVectorCircle.Dot(&nextVertexToCurrVertexVector) < 0 {
		if nextVertexVectorCircle.Magnitude() > circle.Radius {
			return false, nil
		}
		contact.Normal = nextVertexVectorCircle.UnitVector()
		normalScaledRadius := contact.Normal.Scale(circle.Radius)
		contact.Start = contact.B.Postion.Sub(&normalScaledRadius)
		contact.Depth = float32(circle.Radius) - float32(nextVertexVectorCircle.Magnitude())
		depthScaledNormal := contact.Normal.Scale(float64(contact.Depth))
		contact.End = contact.Start.Add(&depthScaledNormal)
		return true, &contact
	} else if distanceCircleEdge < circle.Radius {

		contact.Normal = nextVertex.Sub(&currentVertex).Normal()
		normalScaledRadius := contact.Normal.Scale(circle.Radius)
		contact.Start = contact.B.Postion.Sub(&normalScaledRadius)
		contact.Depth = float32(circle.Radius) - float32(distanceCircleEdge)
		depthScaledNormal := contact.Normal.Scale(float64(contact.Depth))
		contact.End = contact.Start.Add(&depthScaledNormal)
		return true, &contact
	}

	return false, nil

}

func checkPolygonPolygonCollision(a *body, b *body) (bool, *ContactInfo) {
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
