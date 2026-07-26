package game

import (
	"astro/gameerrors"
	"astro/physics"
	"astro/util"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	Img          *ebiten.Image
	Bodies       []*physics.Body
	ScreenWidth  int
	ScreenHeight int
	debug        bool
}

var anchor = physics.Vec2{640, 20}

func (g *Game) Update() error {
	err := g.handleQuit()
	if err != nil {
		return err
	}

	for _, body := range g.Bodies {
		weight := physics.Vec2{0, body.Mass * physics.G * physics.PixelPerMeter}
		wind := physics.Vec2{10 * physics.PixelPerMeter, 0}

		body.Forces = body.Forces.Add(&weight)
		body.Forces = body.Forces.Add(physics.GenerateSpringForce(body, anchor, 10, 546))
		body.Forces = body.Forces.Add(&wind)

		body.Forces = body.Forces.Add(&wind)

		body.CalculateAcceleration()
		body.IntegrateAcceleration(ebiten.ActualTPS())
		body.IntegratePosition(ebiten.ActualTPS())

		// zero forces
		body.Forces[0], body.Forces[1] = 0, 0
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// // sky blue
	screen.Fill(color.RGBA{100, 149, 237, 255})
	vector.FillCircle(screen, float32(anchor[0]), float32(anchor[1]), 10, color.RGBA{255, 0, 0, 255}, true)

	for _, body := range g.Bodies {
		if body.Shape.GetType() == physics.CircleShape {
			util.DrawCircle(body.Shape.(*physics.Circle), &body.Postion, screen, g.debug)

		}

		if body.Shape.GetType() == physics.RectangleShape {
			util.DrawRectangle(body.Shape.(*physics.Rectangle), 0, body.Postion, screen, g.debug)
		}
		vector.StrokeLine(screen, float32(anchor[0]), float32(anchor[1]), float32(body.Postion[0]), float32(body.Postion[1]), 2, color.RGBA{51, 255, 51, 1}, true)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) handleQuit() error {
	terminated := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if terminated {
		return gameerrors.GameExitedError
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		body := &physics.Body{Postion: physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, Mass: 2.2, Shape: physics.NewRectangle(200, 100)}
		g.Bodies = append(g.Bodies, body)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		body := &physics.Body{Postion: physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, Mass: 2.2, Shape: &physics.Circle{Radius: 100}}
		g.Bodies = append(g.Bodies, body)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	return nil
}
