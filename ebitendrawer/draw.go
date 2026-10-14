// package ebitendrawer is a catch all to help me draw geometry using ebiten
package ebitendrawer

import (
	"github.com/the-kada/2dphys/physics"
	"fmt"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

func DrawCircle(circle *physics.Circle, rotation float64, bodyPosition *physics.Vec2, dst *ebiten.Image, debug bool) {
	if debug {
		// x' = a + (x − a)cos θ + (y − b)sin θ
		// y' = b − (x − a)sin θ + (y − b)cos θ
		// find a point on the curcumfrence and rotate ie about the circle center to show rotation
		normalScaledRadius := bodyPosition.Normal().Scale(circle.Radius).Add(bodyPosition)
		radians := rotation * (math.Pi / 180)
		x := bodyPosition[0] + ((normalScaledRadius[0] - bodyPosition[0]) * math.Cos(radians)) + ((normalScaledRadius[1] - bodyPosition[1]) * math.Sin(radians))
		y := bodyPosition[1] - ((normalScaledRadius[0] - bodyPosition[0]) * math.Sin(radians)) + ((normalScaledRadius[1] - bodyPosition[1]) * math.Cos(radians))
		vector.StrokeLine(dst, float32(bodyPosition[0]), float32(bodyPosition[1]), float32(x), float32(y), 3, color.RGBA{1, 255, 51, 1}, false)
		vector.StrokeCircle(dst, float32(bodyPosition[0]), float32(bodyPosition[1]), float32(circle.Radius), 2, color.RGBA{51, 255, 51, 1}, true)
	} else {
		vector.FillCircle(dst, float32(bodyPosition[0]), float32(bodyPosition[1]), float32(circle.Radius), color.White, true)
	}
}

func DrawPolygon(polygon *physics.ConvexPolygon, rotation float64, bodyPosition physics.Vec2, dst *ebiten.Image, debug bool) {
	polygon.UpdateVertices(rotation, bodyPosition)

	var path vector.Path
	path.MoveTo(float32(polygon.WorldVertices[0][0]), float32(polygon.WorldVertices[0][1]))
	for i := 1; i < len(polygon.WorldVertices); i++ {
		path.LineTo(float32(polygon.WorldVertices[i][0]), float32(polygon.WorldVertices[i][1]))
	}

	path.Close()

	var colorScaler ebiten.ColorScale

	if debug {
		colorScaler.ScaleWithColor(color.RGBA{51, 255, 51, 1})
		vector.StrokePath(dst, &path, &vector.StrokeOptions{Width: 2}, &vector.DrawPathOptions{ColorScale: colorScaler})
	} else {
		vector.FillPath(dst, &path, &vector.FillOptions{}, &vector.DrawPathOptions{ColorScale: colorScaler})
	}

}

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
