// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: board_view.go
//
// Project: GoTicTacToe
// Authors:
//   - Alexandre Schmid <alexandre.schmid@edu.heia-fr.ch>
//   - Jeremy Prin <jeremy.prin@edu.heia-fr.ch>
//
// Date: 09 January 2026
//
// Copyright:
//
//	Copyright (c) 2026 HEIA-FR / ISC
//	Haute école d'ingénierie et d'architecture de Fribourg
//	Informatique et Systèmes de Communication
//
// License:
//
//	SPDX-License-Identifier: MIT OR Apache-2.0
//
// Description:
//
//	This file implements BoardView, the widget responsible for rendering the
//	game board grid and drawing player symbols, as well as handling mouse input
//	to translate clicks into grid coordinates.
package ui

import (
	"GoTicTacToe/game"
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Layout constants used by BoardView.
const (
	// cellPaddingRatio defines the padding inside each cell as a fraction of cell size.
	// Example: 0.10 means 10% padding on each side.
	cellPaddingRatio = 0.10

	// half is used for centering operations (equivalent to dividing by 2).
	halfcenter = 0.5

	// two is used for readability when subtracting 2*padding.
	two = 2.0
)

// BoardView is the visual component responsible for rendering the
// Tic-Tac-Toe board and handling user interaction.
type BoardView struct {
	Widget // Embeds Widget: inherits size, position, anchor, LayoutRect(), etc.

	logicBoard  *game.Board      // Reference to the logical board
	OnCellClick func(cx, cy int) // Callback triggered when a cell is clicked

	lastGridW int // Cached grid image width
	lastGridH int // Cached grid image height
}

// NewBoardView creates a new BoardView widget.
//
// Parameters:
// - board: logical board reference (game state)
// - x, y: offset (relative to the widget anchor)
// - size: widget width and height (square board rendering)
// - style: visual styling (background, border, etc.)
// - onClick: callback invoked when a cell is clicked (grid coordinates)
func NewBoardView(
	board *game.Board,
	x, y, size float64,
	style utils.WidgetStyle,
	onClick func(cx, cy int),
) *BoardView {
	view := &BoardView{
		Widget: Widget{
			OffsetX: x,
			OffsetY: y,
			Width:   size,
			Height:  size,
			Anchor:  utils.AnchorCenter,
			Style:   style,
		},
		logicBoard:  board,
		OnCellClick: onClick,
	}

	// Pre-generate the grid once (static background).
	view.ensureGridImage(size, size)

	return view
}

// createGridImage renders the static background grid (background + lines)
// and returns the resulting image.
func (v *BoardView) createGridImage(width, height int) *ebiten.Image {
	img := ebiten.NewImage(width, height)

	// Fill board background.
	img.Fill(v.Style.BackgroundNormal)

	cellWidth := float64(width) / float64(v.logicBoard.Width)
	cellHeight := float64(height) / float64(v.logicBoard.Height)

	thickness := v.Style.BorderWidth
	lineColor := v.Style.BorderColor

	// Draw vertical grid lines.
	for i := 1; i < v.logicBoard.Width; i++ {
		offset := float64(i) * cellWidth

		vert := ebiten.NewImage(int(thickness), height)
		vert.Fill(lineColor)

		opv := &ebiten.DrawImageOptions{}
		opv.GeoM.Translate(offset-thickness*halfcenter, 0)
		img.DrawImage(vert, opv)
	}

	// Draw horizontal grid lines.
	for i := 1; i < v.logicBoard.Height; i++ {
		offset := float64(i) * cellHeight

		hori := ebiten.NewImage(width, int(thickness))
		hori.Fill(lineColor)

		oph := &ebiten.DrawImageOptions{}
		oph.GeoM.Translate(0, offset-thickness*halfcenter)
		img.DrawImage(hori, oph)
	}

	return img
}

// ensureGridImage makes sure the cached grid image exists and matches the given size.
func (v *BoardView) ensureGridImage(width, height float64) {
	w := int(width)
	h := int(height)

	if v.Widget.image == nil || v.lastGridW != w || v.lastGridH != h {
		v.Widget.image = v.createGridImage(w, h)
		v.lastGridW = w
		v.lastGridH = h
	}
}

// Update handles mouse click detection and cell coordinate translation.
func (v *BoardView) Update() {
	rect := v.LayoutRect()
	v.ensureGridImage(rect.Width, rect.Height)

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		vx, vy := rect.X, rect.Y

		// Check if click is inside the board boundaries.
		if float64(mx) >= vx && float64(mx) <= vx+rect.Width &&
			float64(my) >= vy && float64(my) <= vy+rect.Height {

			cellWidth := rect.Width / float64(v.logicBoard.Width)
			cellHeight := rect.Height / float64(v.logicBoard.Height)

			// Convert pixel coordinates -> board grid coordinates.
			gridX := int((float64(mx) - vx) / cellWidth)
			gridY := int((float64(my) - vy) / cellHeight)

			// Trigger callback.
			if v.OnCellClick != nil {
				v.OnCellClick(gridX, gridY)
			}
		}
	}
}

// Draw renders the grid and the player symbols for each occupied cell.
func (v *BoardView) Draw(screen *ebiten.Image) {
	rect := v.LayoutRect()
	v.ensureGridImage(rect.Width, rect.Height)
	vx, vy := rect.X, rect.Y

	// Draw static grid background.
	opGrid := &ebiten.DrawImageOptions{}
	srcW := float64(v.Widget.image.Bounds().Dx())
	srcH := float64(v.Widget.image.Bounds().Dy())
	if srcW != 0 && srcH != 0 {
		opGrid.GeoM.Scale(rect.Width/srcW, rect.Height/srcH)
	}
	opGrid.GeoM.Translate(vx, vy)
	screen.DrawImage(v.Widget.image, opGrid)

	cellWidth := rect.Width / float64(v.logicBoard.Width)
	cellHeight := rect.Height / float64(v.logicBoard.Height)

	// Use the smaller dimension for symbol sizing to maintain aspect ratio.
	cellSize := cellWidth
	if cellHeight < cellSize {
		cellSize = cellHeight
	}

	padding := cellSize * cellPaddingRatio
	usableSize := cellSize - two*padding

	// Draw all symbols.
	for x := 0; x < v.logicBoard.Width; x++ {
		for y := 0; y < v.logicBoard.Height; y++ {
			p := v.logicBoard.Cells[x][y]
			if p == nil || p.Symbol.Image == nil {
				continue
			}

			symbolImg := p.Symbol.Image
			srcWInt, srcHInt := symbolImg.Bounds().Dx(), symbolImg.Bounds().Dy()

			// Determine scaling factor based on the largest symbol dimension.
			maxDim := float64(srcWInt)
			if srcHInt > srcWInt {
				maxDim = float64(srcHInt)
			}
			scale := usableSize / maxDim

			opSym := &ebiten.DrawImageOptions{}
			opSym.Filter = ebiten.FilterLinear // Smooth scaling.

			// Scale first.
			opSym.GeoM.Scale(scale, scale)

			// Position inside the cell, centered.
			symbolW := float64(srcWInt) * scale
			symbolH := float64(srcHInt) * scale
			drawX := vx + float64(x)*cellWidth + (cellWidth-symbolW)*halfcenter
			drawY := vy + float64(y)*cellHeight + (cellHeight-symbolH)*halfcenter

			opSym.GeoM.Translate(drawX, drawY)

			// Tint symbol with the player's color.
			opSym.ColorScale.ScaleWithColor(p.Color)

			screen.DrawImage(symbolImg, opSym)
		}
	}
}
