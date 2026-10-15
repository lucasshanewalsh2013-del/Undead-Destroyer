package gamelogic

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Img *ebiten.Image
	X   float64
	Y   float64
}

func (g *Game) Update() error {
	// Movement speed in pixels per frame
	const speed = 2

	// Check for key presses and move sprite
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.X += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.X -= speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Y += speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Y -= speed
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}

	// Translate to current position
	opts.GeoM.Translate(g.X, g.Y)

	screen.DrawImage(g.Img, opts)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 640, 480
}