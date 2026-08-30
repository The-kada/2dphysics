package game

import (
	"astro/gameerrors"
	"astro/physics"
	"astro/util"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
	Bodies       []*physics.Body
	contacts     []*physics.ContactInfo
	Constraints  []physics.Constraint
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

	for i := 0; i < len(g.Bodies); i++ {
		for j := i + 1; j < len(g.Bodies); j++ {
			if colliding, contact := physics.IsColliding(g.Bodies[i], g.Bodies[j]); colliding {
				g.contacts = append(g.contacts, contact)
			}
		}
	}

	for _, contact := range g.contacts {
		constraint := physics.NewPenetrationConstraint(contact.A, contact.B, contact.Start, contact.End, contact.Normal, 0)
		g.Constraints = append(g.Constraints, constraint)
	}

	for _, body := range g.Bodies {
		weight := physics.Vec2{0, body.Mass * 9.8 * 1 * 50}
		body.Forces = body.Forces.Add(&weight)
	}

	for _, body := range g.Bodies {
		body.IntegrateAcceleration(float64(ebiten.TPS()))
	}

	for _, constraint := range g.Constraints {
		constraint.PreSolve(1.0 / float64(ebiten.TPS()))
	}

	for i := 0; i < 5; i++ {
		for _, constraint := range g.Constraints {
			constraint.Solve()
		}
	}

	for _, body := range g.Bodies {
		body.IntegrateVelocity(float64(ebiten.TPS()))
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, body := range g.Bodies {
		if body.Shape.GetType() == physics.CircleShape {
			util.DrawCircle(body.Shape.(*physics.Circle), body.Rotation, &body.Postion, screen, g.debug)

		}

		if body.Shape.GetType() == physics.RectangleShape {
			util.DrawPolygon(&body.Shape.(*physics.Rectangle).ConvexPolygon, body.Rotation, body.Postion, screen, g.debug)
		}

		if body.Shape.GetType() == physics.ConvexPolygonShape {
			util.DrawPolygon(body.Shape.(*physics.ConvexPolygon), body.Rotation, body.Postion, screen, g.debug)
		}

	}

	if g.debug {
		for _, contact := range g.contacts {
			util.HighlightCollision(screen, contact)
		}
	}

	g.contacts = make([]*physics.ContactInfo, 0, 0)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *Game) handleQuit() error {
	terminated := ebiten.IsKeyPressed(ebiten.KeyEscape)
	if terminated {
		return gameerrors.GameExitedError
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyR) {
	// 	body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 60, 0, 0, 0, physics.NewRectangle(100, 100))
	// 	g.Bodies = append(g.Bodies, body)
	// }

	if inpututil.IsKeyJustPressed(ebiten.KeyC) {
		body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 0, 0, 0, 0, 0, 0, &physics.Circle{Radius: 100})
		g.Bodies = append(g.Bodies, body)
	}

	// if inpututil.IsKeyJustPressed(ebiten.KeyP) {
	// 	body := physics.NewBody(physics.Vec2{float64(g.ScreenWidth) / 2, float64(g.ScreenHeight) / 2}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, physics.Vec2{0, 0}, 2.2, 0, 0, 0, 0, physics.NewPolygon([]physics.Vec2{{0, -50}, {20, 70}, {-70, 70}}))
	// 	g.Bodies = append(g.Bodies, body)
	// }

	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.debug = !g.debug
	}

	return nil
}
