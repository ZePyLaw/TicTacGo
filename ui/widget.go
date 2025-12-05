package ui

import (
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Widget represents a UI element that can be drawn on the screen.
type Widget struct {
	OffsetX float64
	OffsetY float64
	Width   float64
	Height  float64

	image  *ebiten.Image
	Anchor utils.Anchor

	Style utils.WidgetStyle
}

func (w *Widget) AbsPosition() (float64, float64) {
	sw, sh := currentScreenSize()
	if sw == 0 || sh == 0 {
		// Fallback to the window size for environments where UpdateScreenSize wasn't called yet.
		sw, sh = ebiten.WindowSize()
	}
	return utils.ComputeAnchoredPosition(
		w.Anchor,
		w.OffsetX, w.OffsetY,
		w.Width, w.Height,
		sw, sh,
	)
}
