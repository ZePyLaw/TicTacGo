package assets

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

// SymbolType identifies which symbol shape to generate and render.
type SymbolType int

const (
	// CrossSymbol represents a "X" symbol.
	CrossSymbol SymbolType = iota

	// CircleSymbol represents an "O" symbol.
	CircleSymbol

	// TriangleSymbol represents a triangle symbol.
	TriangleSymbol

	// SquareSymbol represents a square symbol.
	SquareSymbol
)

// Procedural symbol rendering constants (in pixels).
const (
	// symbolImageSizePx is the width and height of the generated symbol image.
	symbolImageSizePx = 128

	// symbolLineThicknessPx is the stroke thickness used when drawing symbols.
	symbolLineThicknessPx = 20
)

// Symbol represents a renderable game symbol (shape + generated image).
//
// The Image field contains an Ebiten image ready to be drawn in the UI.
type Symbol struct {
	Type  SymbolType
	Image *ebiten.Image
}

// generateSymbol procedurally draws the requested symbol shape and returns it
// as an Ebiten image.
//
// The symbols are currently drawn with a white stroke. Color and styling could
// later be parameterized to support theming.
func generateSymbol(symbolType SymbolType) *ebiten.Image {
	dc := gg.NewContext(symbolImageSizePx, symbolImageSizePx)
	dc.SetColor(color.White)
	dc.SetLineWidth(symbolLineThicknessPx)

	halfStroke := float64(symbolLineThicknessPx) / 2
	max := float64(symbolImageSizePx)

	switch symbolType {
	case CrossSymbol:
		// Diagonal from top-left to bottom-right
		dc.DrawLine(halfStroke, halfStroke, max-halfStroke, max-halfStroke)
		dc.Stroke()
		// Diagonal from bottom-left to top-right
		dc.DrawLine(halfStroke, max-halfStroke, max-halfStroke, halfStroke)
		dc.Stroke()

	case CircleSymbol:
		radius := (max / 2) - halfStroke
		dc.DrawCircle(max/2, max/2, radius)
		dc.Stroke()

	case TriangleSymbol:
		dc.MoveTo(max/2, halfStroke)
		dc.LineTo(max-halfStroke, max-halfStroke)
		dc.LineTo(halfStroke, max-halfStroke)
		dc.ClosePath()
		dc.Stroke()

	case SquareSymbol:
		// Draw a square inset by half the stroke on all sides.
		side := max - float64(symbolLineThicknessPx)
		dc.DrawRectangle(halfStroke, halfStroke, side, side)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

// NewSymbol creates a new Symbol instance by generating its image procedurally.
//
// The returned Symbol is ready to be drawn by the UI.
func NewSymbol(symbolType SymbolType) *Symbol {
	return &Symbol{
		Type:  symbolType,
		Image: generateSymbol(symbolType),
	}
}
