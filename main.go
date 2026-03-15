package main

import (
	"astro/game"
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
		log.Fatal(err)
	}
}
