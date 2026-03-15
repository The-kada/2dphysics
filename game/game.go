package game

import "github.com/hajimehoshi/ebiten"

type Game struct {
}

func (g *Game) Update(*ebiten.Image) error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 0, 9
}
