package physics

import "math"

type ContactInfo struct {
	A      *Body
	B      *Body
	Start  Vec2
	End    Vec2
	Depth  float32
	Normal Vec2
}

func isColliding(a *Body, b *Body) (bool, []*ContactInfo) {
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

func checkCircleCircleCollision(a *Body, b *Body) (bool, []*ContactInfo) {
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

		return true, []*ContactInfo{{a, b, start, end, float32(depth), normalAlongCenters}}
	}

	return false, nil
}

func checkCirclePolygonCollision(a *Body, b *Body) (bool, []*ContactInfo) {

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
		return true, []*ContactInfo{&contact}
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
		return true, []*ContactInfo{&contact}

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
		return true, []*ContactInfo{&contact}
	} else if distanceCircleEdge < circle.Radius {

		contact.Normal = nextVertex.Sub(&currentVertex).Normal()
		normalScaledRadius := contact.Normal.Scale(circle.Radius)
		contact.Start = contact.B.Postion.Sub(&normalScaledRadius)
		contact.Depth = float32(circle.Radius) - float32(distanceCircleEdge)
		depthScaledNormal := contact.Normal.Scale(float64(contact.Depth))
		contact.End = contact.Start.Add(&depthScaledNormal)
		return true, []*ContactInfo{&contact}
	}

	return false, nil

}

func checkPolygonPolygonCollision(a *Body, b *Body) (bool, []*ContactInfo) {
	var polygonA, polygonB *ConvexPolygon
	var ok bool

	if polygonA, ok = a.Shape.(*ConvexPolygon); !ok {
		polygonA = &a.Shape.(*Rectangle).ConvexPolygon
	}

	if polygonB, ok = b.Shape.(*ConvexPolygon); !ok {
		polygonB = &b.Shape.(*Rectangle).ConvexPolygon
	}

	sepAB, aRefIdx := polygonA.findMinSeperation(polygonB)
	if sepAB >= 0 {
		return false, nil
	}

	sepBA, bRefIdx := polygonB.findMinSeperation(polygonA)
	if sepBA >= 0 {
		return false, nil
	}

	// Choose which polygon is the reference (least-negative separation = more face-on)
	// and which is incident.
	var referenceShape, incidentShape *ConvexPolygon
	var indexReferenceEdge int
	referenceIsA := sepAB > sepBA

	if referenceIsA {
		referenceShape = polygonA
		incidentShape = polygonB
		indexReferenceEdge = aRefIdx
	} else {
		referenceShape = polygonB
		incidentShape = polygonA
		indexReferenceEdge = bRefIdx
	}

	referenceEdge := referenceShape.edgeAtIndex(indexReferenceEdge)
	referenceNormal := referenceEdge.Normal()

	// Find the incident edge: the edge on the incident polygon most anti-parallel
	// to the reference edge's normal.
	incidentIdx := incidentShape.findIncidentEdge(referenceNormal)
	incidentNextIdx := (incidentIdx + 1) % len(incidentShape.WorldVertices)
	v0 := incidentShape.WorldVertices[incidentIdx]
	v1 := incidentShape.WorldVertices[incidentNextIdx]

	/////////////////////////////////////
	// Vclipping
	/////////////////////////////////////
	// Clip the incident edge against every edge of the reference polygon except
	// the reference edge itself.
	contactPoints := []Vec2{v0, v1}
	for i := 0; i < len(referenceShape.WorldVertices); i++ {
		if i == indexReferenceEdge {
			continue
		}
		c0 := referenceShape.WorldVertices[i]
		c1 := referenceShape.WorldVertices[(i+1)%len(referenceShape.WorldVertices)]
		contactPoints = clipSegmentToLine(contactPoints, c0, c1)
		if len(contactPoints) < 2 {
			break
		}
	}

	// The reference vertex is the start of the reference edge; separation is
	// measured from it along the reference normal.
	vref := referenceShape.WorldVertices[indexReferenceEdge]
	var contacts []*ContactInfo

	for _, vclip := range contactPoints {
		separation := vclip.Sub(&vref).Dot(&referenceNormal)
		if separation > 0 {
			continue // point is outside the reference face — discard
		}

		// start = contact point on the incident polygon (inside reference)
		// end   = corresponding point projected onto the reference face
		// Convention matches the rest of the codebase: start on B, end on A,
		// normal from A toward B.
		start := vclip
		pushback := referenceNormal.Scale(-separation)
		end := vclip.Add(&pushback)
		normal := referenceNormal

		if !referenceIsA {
			// Reference shape was B; flip so the normal points from a toward b
			// and start/end follow the expected convention.
			start, end = end, start
			normal = normal.Scale(-1)
		}

		contacts = append(contacts, &ContactInfo{
			A:      a,
			B:      b,
			Start:  start,
			End:    end,
			Depth:  float32(-separation),
			Normal: normal,
		})
	}

	if len(contacts) == 0 {
		return false, nil
	}
	return true, contacts
}
