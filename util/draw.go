// Package Util is a catch all to help me draw geometry without coupling to render system
package util

import (
	"astro/physics"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func DrawRectangle(rect *physics.Rectangle, rotation float64, position physics.Vec2, dst *ebiten.Image, debug bool) {
	rect.UpdateVertices(rotation, position)

	var path vector.Path
	path.MoveTo(float32(rect.WorldVertices[0][0]), float32(rect.WorldVertices[0][1]))
	path.LineTo(float32(rect.WorldVertices[1][0]), float32(rect.WorldVertices[1][1]))
	path.LineTo(float32(rect.WorldVertices[2][0]), float32(rect.WorldVertices[2][1]))
	path.LineTo(float32(rect.WorldVertices[3][0]), float32(rect.WorldVertices[3][1]))
	path.Close()

	var colorScaler ebiten.ColorScale

	if debug {
		colorScaler.ScaleWithColor(color.RGBA{51, 255, 51, 1})
		vector.StrokePath(dst, &path, &vector.StrokeOptions{Width: 2}, &vector.DrawPathOptions{ColorScale: colorScaler})
	} else {
		vector.FillPath(dst, &path, &vector.FillOptions{}, &vector.DrawPathOptions{ColorScale: colorScaler})
	}

}

func DrawCircle(circle *physics.Circle, body *physics.Vec2, dst *ebiten.Image, debug bool) {
	if debug {
		vector.StrokeCircle(dst, float32(body[0]), float32(body[1]), float32(circle.Radius), 2, color.RGBA{51, 255, 51, 1}, true)
	} else {
		vector.FillCircle(dst, float32(body[0]), float32(body[1]), float32(circle.Radius), color.White, true)
	}
}
