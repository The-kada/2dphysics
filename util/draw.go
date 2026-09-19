// package util is a catch all to help me draw geometry without coupling to render system
package util

import (
	"astro/physics"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawCircle(circle *physics.Circle, rotation float64, body *physics.Vec2, dst *ebiten.Image, debug bool) {
	if debug {
		// x' = a + (x − a)cos θ + (y − b)sin θ
		// y' = b − (x − a)sin θ + (y − b)cos θ
		// find a point on the curcumfrence and rotate ie about the circle center to show rotation
		normalScaledRadius := body.Normal().Scale(circle.Radius).Add(body)
		radians := rotation * (math.Pi / 180)
		x := body[0] + ((normalScaledRadius[0] - body[0]) * math.Cos(radians)) + ((normalScaledRadius[1] - body[1]) * math.Sin(radians))
		y := body[1] - ((normalScaledRadius[0] - body[0]) * math.Sin(radians)) + ((normalScaledRadius[1] - body[1]) * math.Cos(radians))
		vector.StrokeLine(dst, float32(body[0]), float32(body[1]), float32(x), float32(y), 3, color.RGBA{1, 255, 51, 1}, false)
		vector.StrokeCircle(dst, float32(body[0]), float32(body[1]), float32(circle.Radius), 2, color.RGBA{51, 255, 51, 1}, true)
	} else {
		vector.FillCircle(dst, float32(body[0]), float32(body[1]), float32(circle.Radius), color.RGBA{125, 255, 255, 0}, true)
	}
}

func DrawPolygon(polygon *physics.ConvexPolygon, rotation float64, position physics.Vec2, dst *ebiten.Image, debug bool) {
	polygon.UpdateVertices(rotation, position)

	var path vector.Path
	path.MoveTo(float32(polygon.WorldVertices[0][0]), float32(polygon.WorldVertices[0][1]))
	for i := 1; i < len(polygon.WorldVertices); i++ {
		path.LineTo(float32(polygon.WorldVertices[i][0]), float32(polygon.WorldVertices[i][1]))
	}

	path.Close()

	var colorScaler ebiten.ColorScale
	colorScaler.ScaleWithColor(color.RGBA{220, 220, 220, 128})

	if debug {
		colorScaler.ScaleWithColor(color.RGBA{51, 255, 51, 1})
		vector.StrokePath(dst, &path, &vector.StrokeOptions{Width: 2}, &vector.DrawPathOptions{ColorScale: colorScaler})
	} else {
		vector.FillPath(dst, &path, &vector.FillOptions{}, &vector.DrawPathOptions{ColorScale: colorScaler})
	}

}
