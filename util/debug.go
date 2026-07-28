package util

import (
	"astro/physics"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func HighlightCollision(screen *ebiten.Image, contact *physics.ContactInfo) {
	vector.FillCircle(screen, float32(contact.Start[0]), float32(contact.Start[1]), 4, color.RGBA{255, 0, 0, 100}, false)
	vector.FillCircle(screen, float32(contact.End[0]), float32(contact.End[1]), 4, color.RGBA{255, 0, 0, 100}, false)
	vector.StrokeLine(screen, float32(contact.Start[0]), float32(contact.Start[1]), float32(contact.End[0]), float32(contact.End[1]), 3, color.RGBA{255, 0, 0, 100}, false)
}
