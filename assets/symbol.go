package assets

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type SymbolType int

const (
	CrossSymbol SymbolType = iota
	CircleSymbol
	TriangleSymbol
	SquareSymbol
)

const (
	imageSize     = 128 // pixels
	lineThickness = 20  // pixels
)

type Symbol struct {
	Type  SymbolType
	Image *ebiten.Image
}

func generateSymbol(symbolType SymbolType) *ebiten.Image {
	dc := gg.NewContext(imageSize, imageSize)
	dc.SetColor(color.White)
	dc.SetLineWidth(lineThickness)

	switch symbolType {
	case CrossSymbol:
		dc.DrawLine(lineThickness/2, lineThickness/2, imageSize-lineThickness/2, imageSize-lineThickness/2)
		dc.Stroke()
		dc.DrawLine(lineThickness/2, imageSize-lineThickness/2, imageSize-lineThickness/2, lineThickness/2)
		dc.Stroke()

	case CircleSymbol:
		dc.DrawCircle(imageSize/2, imageSize/2, imageSize/2-lineThickness/2)
		dc.Stroke()

	case TriangleSymbol:
		dc.MoveTo(imageSize/2, lineThickness/2)
		dc.LineTo(imageSize-lineThickness/2, imageSize-lineThickness/2)
		dc.LineTo(lineThickness/2, imageSize-lineThickness/2)
		dc.ClosePath()
		dc.Stroke()

	case SquareSymbol:
		dc.DrawRectangle(lineThickness/2, lineThickness/2, imageSize-lineThickness, imageSize-lineThickness)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

func NewSymbol(symbolType SymbolType) *Symbol {
	return &Symbol{
		Type:  symbolType,
		Image: generateSymbol(symbolType),
	}
}
