package main

import (
	"astro/game"
	"astro/gameerrors"
	"astro/physics"
	"errors"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 320 * 4
	screenHeight = 288 * 4
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(fmt.Sprintf("Astro %dx%d", screenWidth, screenHeight))

	if err := ebiten.RunGame(NewGame()); err != nil {
		if errors.Is(err, gameerrors.GameExitedError) {
			log.Println("Game exited!")
			return
		}
		log.Fatal(err)
	}
}

func NewGame() *game.Game {
	game := &game.Game{
		Bodies:       make([]*physics.Body, 0),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
		Constraints:  make([]physics.Constraint, 0),
	}

	leftWall := physics.NewBody(physics.Vec2{40, float64(game.ScreenHeight) - float64(game.ScreenHeight)/4}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 0, 0, physics.NewRectangle(80, float64(game.ScreenHeight)/2))
	game.Bodies = append(game.Bodies, leftWall)

	rightWall := physics.NewBody(physics.Vec2{float64(game.ScreenWidth) - 40, float64(game.ScreenHeight) - float64(game.ScreenHeight)/4}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 0, 0, physics.NewRectangle(80, float64(game.ScreenHeight)/2))
	game.Bodies = append(game.Bodies, rightWall)

	floor := physics.NewBody(physics.Vec2{float64(game.ScreenWidth) / 2, float64(game.ScreenHeight) - 40}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, 0, 0, 0, 0, 0, 0, 0, physics.NewRectangle(float64(game.ScreenWidth)-160, 80))
	game.Bodies = append(game.Bodies, floor)

	return game
}
