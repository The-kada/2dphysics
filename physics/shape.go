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
	GetMommentOfInertia() float64
	UpdateVertices(angle float64, position Vec2)
	GetType() ShapeType
}

type Circle struct {
	Radius float64
}

func (c *Circle) GetMommentOfInertia() float64 {
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

func (p *ConvexPolygon) FindMinSeperation(op *ConvexPolygon) (seperation float64, axis Vec2, point Vec2) {
	seperation = -math.MaxFloat64

	for i := 0; i < len(p.WorldVertices); i++ {
		edge := p.edgeAtIndex(i)
		normal := edge.Normal()

		minSep := math.MaxFloat64
		var tempPoint Vec2
		for j := 0; j < len(op.WorldVertices); j++ {
			distance := op.WorldVertices[j].Sub(&p.WorldVertices[i]).Dot(&normal)
			if distance < minSep {
				minSep = distance
				tempPoint = op.WorldVertices[j]
			}
		}

		if minSep > seperation {
			seperation = minSep
			axis = p.edgeAtIndex(i)
			point = tempPoint
		}
	}

	return seperation, axis, point
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

func (p *ConvexPolygon) GetMommentOfInertia() float64 {
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

func (r *Rectangle) GetMommentOfInertia() float64 {
	return (0.083333) * (r.width*r.width + r.height*r.height)
}
