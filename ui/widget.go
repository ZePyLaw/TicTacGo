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

	WidthMode  utils.SizeMode
	HeightMode utils.SizeMode

	parentRect utils.LayoutRect
	hasParent  bool
}

func (w *Widget) AbsPosition() (float64, float64) {
	rect := w.LayoutRect()
	return rect.X, rect.Y
}

// LayoutRect returns the resolved rectangle in absolute coordinates, using either
// the parent bounds or the screen size when the widget lives at the root.
func (w *Widget) LayoutRect() utils.LayoutRect {
	parent := w.parentRect
	if !w.hasParent {
		sw, sh := currentScreenSize()
		if sw == 0 || sh == 0 {
			// Fallback to the window size for environments where UpdateScreenSize wasn't called yet.
			sw, sh = ebiten.WindowSize()
		}
		parent = utils.LayoutRect{
			Width:  float64(sw),
			Height: float64(sh),
		}
	}

	width := w.Width
	height := w.Height

	if w.WidthMode == utils.SizeFill {
		width = parent.Width
	}
	if w.HeightMode == utils.SizeFill {
		height = parent.Height
	}

	x, y := utils.ComputeAnchoredPosition(
		w.Anchor,
		w.OffsetX, w.OffsetY,
		width, height,
		parent.Width, parent.Height,
	)

	return utils.LayoutRect{
		X:      parent.X + x,
		Y:      parent.Y + y,
		Width:  width,
		Height: height,
	}
}

// LayoutSize returns the resolved width and height using the current layout context.
func (w *Widget) LayoutSize() (float64, float64) {
	rect := w.LayoutRect()
	return rect.Width, rect.Height
}

// SetParentBounds assigns the parent rectangle so the widget can resolve its layout relative to it.
func (w *Widget) SetParentBounds(parent utils.LayoutRect) {
	w.parentRect = parent
	w.hasParent = true
}

// ClearParentBounds removes any parent information, causing the widget to lay out at the screen root.
func (w *Widget) ClearParentBounds() {
	w.hasParent = false
}
