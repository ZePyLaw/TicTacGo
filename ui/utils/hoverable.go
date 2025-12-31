package utils

import (
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
