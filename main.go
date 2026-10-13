package main

import (
	"2dphys/ebitendrawer"
	"2dphys/physics"
	"errors"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 320 * 3
	screenHeight = 288 * 3
)

type Game struct {
	*physics.WorldPhysics
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

	g.UpdatePhysics(ebiten.ActualTPS())

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, body := range g.Bodies {
		if body.Shape.GetType() == physics.CircleShape {
			ebitendrawer.DrawCircle(body.Shape.(*physics.Circle), body.Rotation, &body.Postion, screen, g.debug)

		}

		if body.Shape.GetType() == physics.RectangleShape {
			ebitendrawer.DrawPolygon(&body.Shape.(*physics.Rectangle).ConvexPolygon, body.Rotation, body.Postion, screen, g.debug)
		}

		if body.Shape.GetType() == physics.ConvexPolygonShape {
			ebitendrawer.DrawPolygon(body.Shape.(*physics.ConvexPolygon), body.Rotation, body.Postion, screen, g.debug)
		}

	}

	if g.debug {
		for _, contact := range g.Contacts {
			ebitendrawer.HighlightCollision(screen, contact)
		}
	}

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) handleQuit() error {
	terminated := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if terminated {
		return errors.New("err")
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 60, 0, 0, 0, 2, 1.3, physics.NewRectangle(100, 100))
		g.AddBody(body)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 0, 0, 0, 0, 2, 1.3, &physics.Circle{Radius: 100})
		g.AddBody(body)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyP) {
		body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 0, 0, 0, 0, 2, 1.3, physics.NewPolygon([]physics.Vec2{{0, -50}, {20, 70}, {-70, 70}}))
		g.AddBody(body)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	return nil
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(fmt.Sprintf("Astro %dx%d", screenWidth, screenHeight))

	if err := ebiten.RunGame(NewGame()); err != nil {
		if err != nil {
			log.Println("Game exited!")
			return
		}
		log.Fatal(err)
	}
}

func NewGame() *Game {
	game := &Game{

		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		WorldPhysics: &physics.WorldPhysics{},
	}

	leftWall := physics.NewBody(physics.Vec2{40, float64(game.ScreenHeight) - float64(game.ScreenHeight)/4}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 2, 1.3, physics.NewRectangle(80, float64(game.ScreenHeight)/2))
	rightWall := physics.NewBody(physics.Vec2{float64(game.ScreenWidth) - 40, float64(game.ScreenHeight) - float64(game.ScreenHeight)/4}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 2, 1.3, physics.NewRectangle(80, float64(game.ScreenHeight)/2))
	floor := physics.NewBody(physics.Vec2{float64(game.ScreenWidth) / 2, float64(game.ScreenHeight) - 40}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 2, 1.3, physics.NewRectangle(float64(game.ScreenWidth)-160, 80))
	game.WorldPhysics.AddBody(leftWall, rightWall, floor)

	return game
}
