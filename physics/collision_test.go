package physics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func makeCircleBody(x, y, radius float64) *Body {
	return &Body{
		Postion: Vec2{x, y},
		Shape:   &Circle{Radius: radius},
	}
}

func TestIsColliding_CirclesOverlapping(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(8, 0, 5)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
	assert.Equal(t, a, contacts[0].A)
	assert.Equal(t, b, contacts[0].B)
	assert.True(t, contacts[0].Depth > 0)
}

func TestIsColliding_CirclesTouching(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(10, 0, 5)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_CirclesSeparated(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(20, 0, 5)

	colliding, contacts := IsColliding(a, b)

	assert.False(t, colliding)
	assert.Empty(t, contacts)
}

func TestIsColliding_CirclesSamePosition(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(0, 0, 5)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_ContactInfo(t *testing.T) {
	// a at origin with r=3, b at (4,0) with r=3 — overlap of 2 units
	a := makeCircleBody(0, 0, 3)
	b := makeCircleBody(4, 0, 3)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
	assert.True(t, almostEqual(float64(contacts[0].Depth), 2.0, 1e-4))
}

// --- helpers ---

func makePolygonBody(x, y float64, vertices []Vec2) *Body {
	p := NewPolygon(vertices)
	p.UpdateVertices(0, Vec2{x, y})
	return &Body{
		Postion: Vec2{x, y},
		Shape:   p,
	}
}

func makeRectBody(x, y, w, h float64) *Body {
	r := NewRectangle(w, h)
	r.UpdateVertices(0, Vec2{x, y})
	return &Body{
		Postion: Vec2{x, y},
		Shape:   r,
	}
}

// --- circle-polygon ---

func TestIsColliding_CircleInsidePolygon(t *testing.T) {
	// square centered at origin, circle well inside it
	poly := makeRectBody(0, 0, 20, 20)
	circle := makeCircleBody(0, 0, 2)

	colliding, contacts := IsColliding(poly, circle)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_CircleOverlapsPolygonEdge(t *testing.T) {
	// 10x10 square at origin, circle overlapping the right edge
	poly := makeRectBody(0, 0, 10, 10)
	circle := makeCircleBody(7, 0, 3) // edge at x=5, circle center at x=7 with r=3 → 1 unit overlap

	colliding, contacts := IsColliding(poly, circle)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
	assert.True(t, contacts[0].Depth > 0)
}

func TestIsColliding_CircleJustInsidePolygonEdge(t *testing.T) {
	// circle slightly overlapping the right edge (1 unit overlap)
	poly := makeRectBody(0, 0, 10, 10)
	circle := makeCircleBody(7, 0, 3) // edge at x=5, center at x=7 with r=3 → 1 unit overlap

	colliding, contacts := IsColliding(poly, circle)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_CircleSeparatedFromPolygon(t *testing.T) {
	poly := makeRectBody(0, 0, 10, 10)
	circle := makeCircleBody(20, 0, 3)

	colliding, contacts := IsColliding(poly, circle)

	assert.False(t, colliding)
	assert.Empty(t, contacts)
}

func TestIsColliding_CircleNearPolygonCorner_Overlapping(t *testing.T) {
	// circle approaching the top-right corner of the square
	poly := makeRectBody(0, 0, 10, 10) // corner at (5,5)
	circle := makeCircleBody(7, 7, 4)  // distance to corner ≈ 2.83, radius 4 → overlap

	colliding, contacts := IsColliding(poly, circle)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_CircleNearPolygonCorner_Separated(t *testing.T) {
	poly := makeRectBody(0, 0, 10, 10) // corner at (5,5)
	circle := makeCircleBody(9, 9, 2)  // distance to corner ≈ 5.66, radius 2 → no overlap

	colliding, contacts := IsColliding(poly, circle)

	assert.False(t, colliding)
	assert.Empty(t, contacts)
}

// --- polygon-polygon ---

func TestIsColliding_RectanglesOverlapping(t *testing.T) {
	a := makeRectBody(0, 0, 10, 10)
	b := makeRectBody(8, 0, 10, 10) // 2 unit overlap on x-axis

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
	assert.Equal(t, a, contacts[0].A)
	assert.Equal(t, b, contacts[0].B)
	assert.True(t, contacts[0].Depth > 0)
}

func TestIsColliding_RectanglesSeparated(t *testing.T) {
	a := makeRectBody(0, 0, 10, 10)
	b := makeRectBody(20, 0, 10, 10)

	colliding, contacts := IsColliding(a, b)

	assert.False(t, colliding)
	assert.Empty(t, contacts)
}

func TestIsColliding_RectanglesJustOverlapping(t *testing.T) {
	// 1 unit overlap on x-axis
	a := makeRectBody(0, 0, 10, 10)
	b := makeRectBody(9, 0, 10, 10)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_RectanglesSamePosition(t *testing.T) {
	a := makeRectBody(0, 0, 10, 10)
	b := makeRectBody(0, 0, 10, 10)

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
}

func TestIsColliding_PolygonPolygonOverlapping(t *testing.T) {
	// two triangles overlapping
	a := makePolygonBody(0, 0, []Vec2{{-5, -5}, {5, -5}, {0, 5}})
	b := makePolygonBody(4, 0, []Vec2{{-5, -5}, {5, -5}, {0, 5}})

	colliding, contacts := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotEmpty(t, contacts)
	assert.True(t, contacts[0].Depth > 0)
}

func TestIsColliding_PolygonPolygonSeparated(t *testing.T) {
	a := makePolygonBody(0, 0, []Vec2{{-5, -5}, {5, -5}, {0, 5}})
	b := makePolygonBody(30, 0, []Vec2{{-5, -5}, {5, -5}, {0, 5}})

	colliding, contacts := IsColliding(a, b)

	assert.False(t, colliding)
	assert.Empty(t, contacts)
}
