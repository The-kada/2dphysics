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

	for i := 0; i < 8; i++ {
		mass := 0.0
		if i != 0 {
			mass = 1
		}
		body := physics.NewBody(physics.Vec2{float64(game.ScreenWidth)/2 - float64(i*40), 100}, physics.Vec2{}, physics.Vec2{}, physics.Vec2{}, mass, 0, 0, 0, 0, physics.NewRectangle(30, 30))
		game.Bodies = append(game.Bodies, body)
	}

	for i := 0; i < len(game.Bodies)-1; i++ {
		constraint := physics.NewJointConstraint(game.Bodies[i], game.Bodies[i+1], game.Bodies[i].Postion)
		game.Constraints = append(game.Constraints, constraint)
	}

	return game
}
