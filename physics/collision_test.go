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

	colliding, contact := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotNil(t, contact)
	assert.Equal(t, a, contact.A)
	assert.Equal(t, b, contact.B)
	assert.True(t, contact.Depth > 0)
}

func TestIsColliding_CirclesTouching(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(10, 0, 5)

	colliding, contact := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotNil(t, contact)
}

func TestIsColliding_CirclesSeparated(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(20, 0, 5)

	colliding, contact := IsColliding(a, b)

	assert.False(t, colliding)
	assert.Nil(t, contact)
}

func TestIsColliding_CirclesSamePosition(t *testing.T) {
	a := makeCircleBody(0, 0, 5)
	b := makeCircleBody(0, 0, 5)

	colliding, contact := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotNil(t, contact)
}

func TestIsColliding_ContactInfo(t *testing.T) {
	// a at origin with r=3, b at (4,0) with r=3 — overlap of 2 units
	a := makeCircleBody(0, 0, 3)
	b := makeCircleBody(4, 0, 3)

	colliding, contact := IsColliding(a, b)

	assert.True(t, colliding)
	assert.NotNil(t, contact)
	assert.True(t, almostEqual(float64(contact.Depth), 2.0, 1e-4))
}
