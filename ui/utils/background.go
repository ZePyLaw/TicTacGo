package utils

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// CreateGradientBackground creates a new image with a vertical gradient.
func CreateGradientBackground(width, height int, top, bottom color.Color) *ebiten.Image {
	img := ebiten.NewImage(width, height)
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		t := float64(y-bounds.Min.Y) / float64(bounds.Dy())
		c := LerpColor(top, bottom, t)
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Set(x, y, c)
		}
	}
	return img
}
