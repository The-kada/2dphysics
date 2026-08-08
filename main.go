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
	}

	//body := &physics.Body{Postion: physics.Vec2{0, 0}, Mass: 2.2, Shape: physics.NewPolygon([]physics.Vec2{{0, -50}, {20, 70}, {-70, 70}})}
	// body := &physics.Body{Postion: physics.Vec2{0, 0}, Mass: 2.2, Shape: physics.NewRectangle(100, 100)}
	body := &physics.Body{Postion: physics.Vec2{0, 0}, Mass: 2.2, Shape: &physics.Circle{Radius: 30}}
	game.Bodies = append(game.Bodies, body)

	return game
}
