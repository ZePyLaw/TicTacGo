package main

import (
	"image/color"

	"github.com/fogleman/gg"
	"github.com/hajimehoshi/ebiten/v2"
)

type SymbolType int

const (
	CROSS SymbolType = iota
	CIRCLE
)

const (
	imageSize = 128 // pixels
	lineThickness = 20 // pixels
)

type Symbol struct {
	Type SymbolType
	image *ebiten.Image
}

func generateImage(symbolType SymbolType) *ebiten.Image {
	dc := gg.NewContext(imageSize, imageSize)
	dc.SetColor(color.White)
	dc.SetLineWidth(lineThickness)

	switch symbolType {
	case CROSS:
		dc.DrawLine(lineThickness/2, lineThickness/2, imageSize-lineThickness/2, imageSize-lineThickness/2)
		dc.Stroke()
		dc.DrawLine(lineThickness/2, imageSize-lineThickness/2, imageSize-lineThickness/2, lineThickness/2)
		dc.Stroke()

	case CIRCLE:
		dc.DrawCircle(imageSize/2, imageSize/2, imageSize/2 - lineThickness/2)
		dc.Stroke()
	}

	return ebiten.NewImageFromImage(dc.Image())
}

func newSymbol(symbolType SymbolType) *Symbol {
	return &Symbol{
		Type:  symbolType,
		image: generateImage(symbolType),
	}
}
