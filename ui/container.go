// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: container.go
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
//	This file implements Container, a layout widget that positions and draws
//	child UI elements relative to its own anchored rectangle. It supports
//	optional padding (inner content rect).
package ui

import (
	"GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

// Numeric constants used for layout and scaling computations.
const (
	zero = 0.0
)

// Element is the minimal interface required by Container to position and render children.
//
// A child element receives the container "inner rect" through SetParentBounds,
// then it is updated and drawn every frame.
type Element interface {
	SetParentBounds(utils.LayoutRect)
	Update()
	Draw(screen *ebiten.Image)
}

// Container is a layout helper that anchors itself like any other widget
// and lets you place child elements relative to its local space.
//
// Container can apply Padding to define an inner rectangle where children are laid out.
// The background is pre-rendered into an image for efficient rendering.
type Container struct {
	Widget
	Padding  utils.Insets
	children []Element
}

// NewContainer creates a rectangular container with a pre-rendered background.
//
// Parameters:
// - offsetX, offsetY: offset relative to the chosen anchor
// - width, height: container size
// - anchor: anchor type (center, top-left, etc.)
// - radius: background corner radius
// - style: widget style (background color, etc.)
func NewContainer(
	offsetX, offsetY float64,
	width, height float64,
	anchor utils.Anchor,
	radius float64,
	style utils.WidgetStyle,
) *Container {
	bg := utils.CreateRoundedRect(int(width), int(height), radius, style.BackgroundNormal)

	return &Container{
		Widget: Widget{
			OffsetX: offsetX,
			OffsetY: offsetY,
			Width:   width,
			Height:  height,
			Anchor:  anchor,
			image:   bg,
			Style:   style,
		},
		Padding: utils.Insets{},
	}
}

// AddChild appends a child element to the container.
func (c *Container) AddChild(child Element) {
	c.children = append(c.children, child)
}

// Update forwards input to children using the container's anchored position.
func (c *Container) Update() {
	_, inner := c.resolveRects()
	c.updateChildren(inner)
}

// Draw renders the container background and all children relative to it.
func (c *Container) Draw(screen *ebiten.Image) {
	outer, inner := c.resolveRects()

	// Draw background (if present).
	if c.image != nil {
		op := &ebiten.DrawImageOptions{}

		srcW := float64(c.image.Bounds().Dx())
		srcH := float64(c.image.Bounds().Dy())

		scaleX := one
		scaleY := one

		// Protect against division by zero (should not happen, but safe).
		if srcW != zero {
			scaleX = outer.Width / srcW
		}
		if srcH != zero {
			scaleY = outer.Height / srcH
		}

		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(outer.X, outer.Y)
		screen.DrawImage(c.image, op)
	}

	// Draw children inside the inner rect.
	for _, child := range c.children {
		child.SetParentBounds(inner)
		child.Draw(screen)
	}
}

// resolveRects returns the outer (container) rect and the inner (content) rect,
// applying padding and clamping inner dimensions to a minimum of zero.
func (c *Container) resolveRects() (utils.LayoutRect, utils.LayoutRect) {
	outer := c.LayoutRect()

	inner := utils.LayoutRect{
		X:      outer.X + c.Padding.Left,
		Y:      outer.Y + c.Padding.Top,
		Width:  outer.Width - c.Padding.Horizontal(),
		Height: outer.Height - c.Padding.Vertical(),
	}

	if inner.Width < zero {
		inner.Width = zero
	}
	if inner.Height < zero {
		inner.Height = zero
	}

	return outer, inner
}

// updateChildren updates children using the provided inner bounds.
func (c *Container) updateChildren(inner utils.LayoutRect) {
	for _, child := range c.children {
		child.SetParentBounds(inner)
		child.Update()
	}
}
