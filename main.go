package main

import (
	"astro/game"
	"astro/gameerrors"
	"astro/physics"
	"errors"
	"fmt"
	"image/color"
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

	img := ebiten.NewImage(500, 500)
	img.Fill(color.White)

	game := &game.Game{
		Img:          img,
		Bodies:       make([]*physics.Body, 0),
		ScreenWidth:  screenWidth,
		ScreenHeight: screenHeight,
	}
	if err := ebiten.RunGame(game); err != nil {
		if errors.Is(err, gameerrors.GameExitedError) {
			log.Println("Game exited!")
			return
		}
		log.Fatal(err)
	}
}
