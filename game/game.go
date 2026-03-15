package game

import (
	"astro/gameerrors"

	"github.com/hajimehoshi/ebiten"
)

type Game struct {
}

func (g *Game) Update(*ebiten.Image) error {
	return handleQuit()
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func handleQuit() error {
	terminated := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if terminated {
		return gameerrors.GameExitedError
	}
	return nil
}
