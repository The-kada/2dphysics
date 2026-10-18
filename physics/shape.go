package physics

import (
	"math"
)

type ShapeType string

const (
	CircleShape        ShapeType = "circle"
	RectangleShape     ShapeType = "rectangle"
	ConvexPolygonShape ShapeType = "ConvexPolygon"
)

type shape interface {
	getMommentOfInertia() float64
	UpdateVertices(angle float64, position Vec2)
	GetType() ShapeType
}

type Circle struct {
	Radius float64
}

func (c *Circle) getMommentOfInertia() float64 {
	return 0.5 * c.Radius * c.Radius
}

func (c *Circle) UpdateVertices(angle float64, position Vec2) {
	// no-op
}

func (c *Circle) GetType() ShapeType {
	return CircleShape
}

type ConvexPolygon struct {
	localVertices []Vec2
	WorldVertices []Vec2
}

func (p *ConvexPolygon) findMinSeperation(op *ConvexPolygon) (separation float64, referenceEdgeIndex int) {
	separation = -math.MaxFloat64

	for i := 0; i < len(p.WorldVertices); i++ {
		edge := p.edgeAtIndex(i)
		normal := edge.Normal()

		minSep := math.MaxFloat64
		for j := 0; j < len(op.WorldVertices); j++ {
			distance := op.WorldVertices[j].Sub(&p.WorldVertices[i]).Dot(&normal)
			if distance < minSep {
				minSep = distance
			}
		}

		if minSep > separation {
			separation = minSep
			referenceEdgeIndex = i
		}
	}

	return separation, referenceEdgeIndex
}

// findIncidentEdge returns the index of the edge on p whose outward normal is
// most anti-parallel to referenceNormal (i.e. the edge most facing the reference edge).
func (p *ConvexPolygon) findIncidentEdge(referenceNormal Vec2) int {
	incidentIndex := 0
	minProj := math.MaxFloat64
	for i := 0; i < len(p.WorldVertices); i++ {
		edgeNormal := p.edgeAtIndex(i).Normal()
		proj := edgeNormal.Dot(&referenceNormal)
		if proj < minProj {
			minProj = proj
			incidentIndex = i
		}
	}
	return incidentIndex
}

// clipSegmentToLine performs one step of Sutherland-Hodgman clipping:
// it clips the input polygon/segment against the half-plane defined by edge c0→c1,
// keeping points on the interior side (dot with outward normal <= 0).
func clipSegmentToLine(input []Vec2, c0, c1 Vec2) []Vec2 {
	if len(input) == 0 {
		return nil
	}
	output := make([]Vec2, 0, len(input))
	c0c1 := c1.Sub(&c0)
	n := c0c1.Normal() // outward normal of the clipping edge

	for i := 0; i < len(input); i++ {
		curr := input[i]
		next := input[(i+1)%len(input)]

		dCurr := curr.Sub(&c0).Dot(&n)
		dNext := next.Sub(&c0).Dot(&n)

		if dCurr <= 0 {
			output = append(output, curr)
		}
		// If curr and next are on opposite sides, add the intersection point.
		if (dCurr < 0 && dNext > 0) || (dCurr > 0 && dNext < 0) {
			t := dCurr / (dCurr - dNext)
			diff := next.Sub(&curr)
			scaled := diff.Scale(t)
			output = append(output, curr.Add(&scaled))
		}
	}
	return output
}

func (p *ConvexPolygon) edgeAtIndex(index int) Vec2 {
	currVertex := p.WorldVertices[index]
	nextVertex := p.WorldVertices[(index+1)%len(p.WorldVertices)]
	return nextVertex.Sub(&currVertex)
}

func NewPolygon(vertices []Vec2) *ConvexPolygon {
	p := &ConvexPolygon{}

	p.WorldVertices = append(p.WorldVertices, vertices...)
	p.localVertices = append(p.localVertices, vertices...)

	return p
}

func (p *ConvexPolygon) getMommentOfInertia() float64 {
	// todo
	return 5000
}

func (p *ConvexPolygon) UpdateVertices(angle float64, position Vec2) {
	for i := 0; i < len(p.localVertices); i++ {
		p.WorldVertices[i] = p.localVertices[i].Rotate(angle)
		p.WorldVertices[i] = p.WorldVertices[i].Add(&position)
	}
}

func (p *ConvexPolygon) GetType() ShapeType {
	return ConvexPolygonShape
}

type Rectangle struct {
	width  float64
	height float64
	ConvexPolygon
}

func NewRectangle(width, height float64) *Rectangle {
	rectangle := &Rectangle{width: width, height: height}
	rectangle.localVertices = append(rectangle.localVertices,
		Vec2{-width / 2.0, -height / 2.0},
		Vec2{width / 2.0, -height / 2.0},
		Vec2{width / 2.0, height / 2.0},
		Vec2{-width / 2.0, height / 2.0},
	)

	rectangle.WorldVertices = append(rectangle.WorldVertices,
		Vec2{-width / 2.0, -height / 2.0},
		Vec2{width / 2.0, -height / 2.0},
		Vec2{width / 2.0, height / 2.0},
		Vec2{-width / 2.0, height / 2.0},
	)

	return rectangle
}

func (rect *Rectangle) GetType() ShapeType {
	return RectangleShape
}

func (r *Rectangle) getMommentOfInertia() float64 {
	return (0.083333) * (r.width*r.width + r.height*r.height)
}
