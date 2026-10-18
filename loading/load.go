package loading

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/lucaswalsh11/UndeadDestroyer/gamelogic"
)

func loadPNG(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open %s: %w", path, err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode %s: %w", path, err)
	}

	return img, nil
}

func OpenGame() {
	// 1. Load player sprite.
	charImg, err := loadPNG("exportedSprites/Character Model.png")
	if err != nil {
		panic(err)
	}

	// 2. Load bullet sprite.
	bulletImg, err := loadPNG("exportedSprites/Bullet.png")
	if err != nil {
		panic(err)
	}

	// 3. Load regular zombie sprite.
	zombieImg, err := loadPNG("exportedSprites/Zombie.png")
	if err != nil {
		panic(err)
	}

	// 4. Load the separate big zombie sprite.
	bigZombieImg, err := loadPNG("exportedSprites/BigZombie.png")
	if err != nil {
		panic(err)
	}

	// 5. Convert the images to Ebiten images.
	player := ebiten.NewImageFromImage(charImg)
	bullet := ebiten.NewImageFromImage(bulletImg)
	zombie := ebiten.NewImageFromImage(zombieImg)
	bigZombie := ebiten.NewImageFromImage(bigZombieImg)

	// 6. Create the game with all four sprites.
	game := gamelogic.NewGame(
		player,
		bullet,
		zombie,
		bigZombie,
	)

	// 7. Start the game.
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Undead Destroyer")

	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}