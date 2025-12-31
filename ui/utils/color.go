package utils

import "image/color"

func lerp(a, b, t float64) float64 {
	return a + (b-a)*t
}

func LerpColor(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()

	return color.RGBA{
		R: uint8(lerp(float64(r1>>8), float64(r2>>8), t)),
		G: uint8(lerp(float64(g1>>8), float64(g2>>8), t)),
		B: uint8(lerp(float64(b1>>8), float64(b2>>8), t)),
		A: uint8(lerp(float64(a1>>8), float64(a2>>8), t)),
	}
}
