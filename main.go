package main

import (
	"astro/game"
	"astro/gameerrors"
	"errors"
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
)

const (
	screenWidth  = 320
	screenHeight = 288
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle(fmt.Sprintf("Astro %dx%d", screenWidth, screenHeight))

	game := &game.Game{}
	if err := ebiten.RunGame(game); err != nil {
		if errors.Is(err, gameerrors.GameExitedError) {
			log.Println("Game exited!")
			return
		}
		log.Fatal(err)
	}
}
