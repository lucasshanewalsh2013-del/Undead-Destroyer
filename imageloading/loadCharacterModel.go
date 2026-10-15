package loading

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func loadCharacterModel(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	img, err := png.Decode(file)
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()
	if width != 32 || height != 32 {
		return nil, fmt.Errorf("image dimensions are not 32x32: got %dx%d", width, height)
	}
	return img, nil
}