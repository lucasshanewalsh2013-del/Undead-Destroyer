package loading

import (
	"fmt"
	"image/png"
	"os"

	"github.com/lucaswalsh11/UndeadDestroyer/gamelogic"
	"github.com/hajimehoshi/ebiten/v2"
)


func OpenCharacterModel(){

	file, err := os.Open("exportedSprites/Character Model.png")
		if err != nil {
		return
	}
	defer file.Close()

	scrimg, err := png.Decode(file)
	if err != nil {
		return
	}

	// Get sprite dimensions for centering
	width := scrimg.Bounds().Dx()
	height := scrimg.Bounds().Dy()

	game := &gamelogic.Game{
		Img: ebiten.NewImageFromImage(scrimg),
		// Start centered on a 640x480 screen
		X: float64(640-width) / 2,
		Y: float64(480-height) / 2,
	}

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Undead Destroyer")

	if err := ebiten.RunGame(game); err != nil {
		panic(fmt.Errorf("game failed: %w", err))
	}
}