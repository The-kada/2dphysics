package util

import (
	"astro/physics"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func HighlightCollision(screen *ebiten.Image, contact *physics.ContactInfo) {
	vector.FillCircle(screen, float32(contact.Start[0]), float32(contact.Start[1]), 4, color.RGBA{255, 0, 0, 100}, false)
	vector.FillCircle(screen, float32(contact.End[0]), float32(contact.End[1]), 4, color.RGBA{255, 0, 0, 100}, false)
	vector.StrokeLine(screen, float32(contact.Start[0]), float32(contact.Start[1]), float32(contact.End[0]), float32(contact.End[1]), 3, color.RGBA{255, 0, 0, 100}, false)

	op := &text.DrawOptions{}
	op.GeoM.Translate(contact.A.Postion[0]-10, contact.A.Postion[1]-10)
	text.Draw(screen, fmt.Sprint("A", " ", contact.Depth), text.NewGoXFace(basicfont.Face7x13), op)

	op = &text.DrawOptions{}
	op.GeoM.Translate(contact.B.Postion[0]-10, contact.B.Postion[1]-10)
	text.Draw(screen, fmt.Sprint("B", " ", contact.Depth), text.NewGoXFace(basicfont.Face7x13), op)
}
