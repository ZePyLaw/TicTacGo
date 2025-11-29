package utils

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Hoverable struct {
	hover float64
	mode  HoverMode
}

func NewHoverable(mode HoverMode) *Hoverable {
	return &Hoverable{
		hover: 0,
		mode:  mode,
	}
}

func (h *Hoverable) Update(x, y, width, height float64) {
	mx, my := ebiten.CursorPosition()
	hover := float64(mx) >= x &&
		float64(mx) <= x+width &&
		float64(my) >= y &&
		float64(my) <= y+height

	if hover {
		h.hover += 0.1
		if h.hover > 1 {
			h.hover = 1
		}
	} else {
		h.hover -= 0.1
		if h.hover < 0 {
			h.hover = 0
		}
	}
}

func (h *Hoverable) IsHovered() bool {
	return h.hover > 0
}

func (h *Hoverable) ApplyHoverColor(op *ebiten.DrawImageOptions, style WidgetStyle) {
	switch h.mode {
	case HoverFade:
		op.ColorScale.ScaleWithColor(style.BackgroundNormal)
		op.ColorScale.ScaleAlpha(float32(h.hover))
	case HoverColorLerp:
		col := LerpColor(style.BackgroundNormal, style.BackgroundHover, h.hover)
		op.ColorScale.ScaleWithColor(col)
	case HoverSolid:
		if h.hover >= 1 {
			op.ColorScale.ScaleWithColor(style.BackgroundHover)
		} else {
			op.ColorScale.ScaleWithColor(style.BackgroundNormal)
		}
	}
}

type HoverMode int

const (
	HoverFade HoverMode = iota
	HoverColorLerp
	HoverSolid
)

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
