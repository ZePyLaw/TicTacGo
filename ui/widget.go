// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: widget.go
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
//	This file defines Widget, the base building block of the UI system.
//	Widgets handle layout, anchoring, sizing, and parent-child relationships.
//	Higher-level components (buttons, containers, views) embed Widget to
//	inherit this behavior.
package ui

import (
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Widget represents a UI element that can be positioned and sized on screen.
//
// It provides:
//   - anchoring (center, top-left, etc.)
//   - optional parent-relative layout
//   - fixed or fill sizing modes
//
// Widget is intended to be embedded by higher-level UI components.
type Widget struct {
	OffsetX float64 // Horizontal offset relative to anchor
	OffsetY float64 // Vertical offset relative to anchor
	Width   float64 // Desired width
	Height  float64 // Desired height

	image  *ebiten.Image // Optional pre-rendered background image
	Anchor utils.Anchor  // Anchor used to position the widget

	Style utils.WidgetStyle // Visual style (colors, borders, etc.)

	WidthMode  utils.SizeMode // Width sizing mode (fixed or fill)
	HeightMode utils.SizeMode // Height sizing mode (fixed or fill)

	parentRect utils.LayoutRect // Parent layout rectangle
	hasParent  bool             // Whether the widget has a parent container
}

// AbsPosition returns the absolute top-left position of the widget
// in screen coordinates.
func (w *Widget) AbsPosition() (float64, float64) {
	rect := w.LayoutRect()
	return rect.X, rect.Y
}

// LayoutRect returns the resolved layout rectangle in absolute coordinates.
//
// If the widget has a parent, layout is resolved relative to the parent bounds.
// Otherwise, the widget is laid out relative to the full screen size.
//
// If the screen size is not yet known (e.g., early startup), this function
// falls back to Ebiten's WindowSize.
func (w *Widget) LayoutRect() utils.LayoutRect {
	parent := w.parentRect

	if !w.hasParent {
		sw, sh := currentScreenSize()
		if sw == 0 || sh == 0 {
			// Fallback to the window size if UpdateScreenSize was not called yet.
			sw, sh = ebiten.WindowSize()
		}
		parent = utils.LayoutRect{
			Width:  float64(sw),
			Height: float64(sh),
		}
	}

	width := w.Width
	height := w.Height

	// Resolve size modes.
	if w.WidthMode == utils.SizeFill {
		width = parent.Width
	}
	if w.HeightMode == utils.SizeFill {
		height = parent.Height
	}

	// Compute anchored position inside parent.
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

// SetParentBounds assigns a parent layout rectangle.
//
// This causes the widget to resolve its layout relative to the given parent.
func (w *Widget) SetParentBounds(parent utils.LayoutRect) {
	w.parentRect = parent
	w.hasParent = true
}

// ClearParentBounds removes any parent association.
//
// After calling this, the widget will be laid out relative to the screen root.
func (w *Widget) ClearParentBounds() {
	w.hasParent = false
}
