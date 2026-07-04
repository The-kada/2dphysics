package physics

type ShapeType string

const (
	CircleShape ShapeType = "circle"
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

}

func (c *Circle) GetType() ShapeType {
	return CircleShape
}
