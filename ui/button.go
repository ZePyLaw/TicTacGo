// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: button.go
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
//	This file implements a clickable Button widget with hover effects and
//	centered label rendering.
package ui

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// Button rendering and layout constants.
const (
	// one is used for default scaling (no scale).
	one = 1.0

	// half is used for centering computations (equivalent to dividing by 2).
	half = 0.5

	// defaultLineSpacingPx is the line spacing used when drawing the label.
	defaultLineSpacingPx = 20
)

// Button represents a UI button with a label, hover animation, and click handler.
//
// It embeds Widget, so all Widget fields and methods are promoted and accessible
// directly. It also embeds utils.Hoverable to implement hover animations/effects.
type Button struct {
	Widget                  // Embedded base widget
	Label            string // Text displayed on the button
	onClick          func() // Callback executed when the button is clicked
	*utils.Hoverable        // Embedded hoverable component
}

// NewButton creates and returns a new Button instance.
//
// The background is pre-rendered as a rounded rectangle for efficiency.
func NewButton(
	label string,
	offsetX, offsetY float64,
	anchor utils.Anchor,
	width, height float64,
	radius float64,
	style utils.WidgetStyle,
	onClick func(),
) *Button {
	// Pre-render the button background (rounded rectangle).
	bg := utils.CreateRoundedRect(
		int(width), int(height),
		radius, style.BackgroundNormal,
	)

	return &Button{
		Widget: Widget{
			OffsetX: offsetX,
			OffsetY: offsetY,
			Width:   width,
			Height:  height,
			Anchor:  anchor,
			image:   bg, // internal widget image
			Style:   style,
		},
		Label:     label,
		onClick:   onClick,
		Hoverable: utils.NewHoverable(style.HoverMode),
	}
}

// Update handles hover detection, hover animation, and click events.
func (b *Button) Update() {
	rect := b.LayoutRect()
	b.UpdateAt(rect.X, rect.Y)
}

// UpdateAt runs hover + click handling at an explicit screen position.
//
// Useful when the button is drawn inside a translated container.
func (b *Button) UpdateAt(x, y float64) {
	width, height := b.LayoutSize()
	b.Hoverable.Update(x, y, width, height)

	// Handle click events.
	if b.IsHovered() && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		if b.onClick != nil {
			b.onClick()
		}
	}
}

// Draw renders the button background and its centered label.
func (b *Button) Draw(screen *ebiten.Image) {
	rect := b.LayoutRect()
	b.DrawAt(screen, rect.X, rect.Y)
}

// SetStyle updates the button style and rebuilds its background.
func (b *Button) SetStyle(style utils.WidgetStyle, radius float64) {
	b.Style = style
	b.image = utils.CreateRoundedRect(int(b.Width), int(b.Height), radius, style.BackgroundNormal)
}

// DrawAt renders the button at an explicit position without recomputing layout.
//
// Useful when nesting the button inside a container that manages its own origin.
func (b *Button) DrawAt(screen *ebiten.Image, x, y float64) {
	width, height := b.LayoutSize()
	srcW := float64(b.image.Bounds().Dx())
	srcH := float64(b.image.Bounds().Dy())

	scaleX := one
	scaleY := one
	if srcW != 0 {
		scaleX = width / srcW
	}
	if srcH != 0 {
		scaleY = height / srcH
	}

	// Draw background.
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scaleX, scaleY)
	op.GeoM.Translate(x, y)

	// Apply interactive hover color effect.
	b.ApplyHoverColor(op, b.Style)

	screen.DrawImage(b.image, op)

	// Draw text.
	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter
	opts.LineSpacing = defaultLineSpacingPx

	// Center the text inside the button.
	opts.GeoM.Translate(x+width*half, y+height*half)
	opts.ColorScale.ScaleWithColor(b.Style.TextColor)

	text.Draw(screen, b.Label, assets.NormalFont, opts)
}
