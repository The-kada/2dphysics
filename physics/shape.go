package physics

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

func (rect *Rectangle) Draw(dst *ebiten.Image) {
	var path vector.Path
	path.MoveTo(float32(rect.WorldVertices[0][0]), float32(rect.WorldVertices[0][1]))
	path.LineTo(float32(rect.WorldVertices[1][0]), float32(rect.WorldVertices[1][1]))
	path.LineTo(float32(rect.WorldVertices[2][0]), float32(rect.WorldVertices[2][1]))
	path.LineTo(float32(rect.WorldVertices[3][0]), float32(rect.WorldVertices[3][1]))
	path.Close()

	color := ebiten.ColorScale{}
	color.Scale(0.2, 0.6, 1.0, 1.0)
	vector.FillPath(dst, &path, &vector.FillOptions{}, &vector.DrawPathOptions{ColorScale: color})
}

func (rect *Rectangle) GetType() ShapeType {
	return RectangleShape
}

func (r *Rectangle) GetMommentOfInertia() float64 {
	return (0.083333) * (r.width*r.width + r.height*r.height)
}
