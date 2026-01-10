// Package ui contains reusable UI widgets and views rendered with Ebiten.
//
// File: score_view.go
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
//	This file implements ScoreView, a widget displaying the current scores and
//	symbols for any number of players. Non-active players can be visually dimmed
//	while the game is running.
package ui

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/game"
	"GoTicTacToe/ui/utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ScoreView layout constants.
const (
	// Panel appearance.
	scorePanelCornerRadiusPx = 10
	scorePanelOffsetY        = 20

	// Per-player zone layout.
	zonePaddingPx         = 12.0
	zoneMinIconSizePx     = 16.0
	zoneIconBottomPadding = 16.0
	zoneNameYOffsetPx     = 2.0
	zoneNamePaddingRatio  = 0.5 // padding/2

	// Symbol placement.
	zoneIconExtraTopShiftPx = 12.0

	// Visual effect when it's not the player's turn.
	nonActiveAlphaScale = 0.5
)

// ScoreView displays player icons and scores for any number of players.
type ScoreView struct {
	Widget
	gameRef *game.Game
}

// NewScoreView creates a score panel widget.
//
// The widget is anchored at the top-center by default and slightly offset downwards.
func NewScoreView(g *game.Game, width, height float64, style utils.WidgetStyle) *ScoreView {
	bg := utils.CreateRoundedRect(int(width), int(height), scorePanelCornerRadiusPx, style.BackgroundNormal)

	return &ScoreView{
		Widget: Widget{
			Width:   width,
			Height:  height,
			image:   bg,
			Anchor:  utils.AnchorTopCenter,
			OffsetY: scorePanelOffsetY,
			Style:   style,
		},
		gameRef: g,
	}
}

// Draw renders the score panel and each player's information.
func (sv *ScoreView) Draw(screen *ebiten.Image) {
	rect := sv.LayoutRect()
	x, y := rect.X, rect.Y

	// Draw background rectangle.
	op := &ebiten.DrawImageOptions{}
	srcW := float64(sv.image.Bounds().Dx())
	srcH := float64(sv.image.Bounds().Dy())
	if srcW != 0 && srcH != 0 {
		op.GeoM.Scale(rect.Width/srcW, rect.Height/srcH)
	}
	op.GeoM.Translate(x, y)
	screen.DrawImage(sv.image, op)

	playerCount := len(sv.gameRef.Players)
	if playerCount == 0 {
		return
	}

	// Split the width into equal zones for each player.
	zoneWidth := rect.Width / float64(playerCount)

	for i, p := range sv.gameRef.Players {
		zoneX := x + float64(i)*zoneWidth
		sv.drawPlayerZone(screen, p, zoneX, y, zoneWidth, rect.Height)
	}
}

// drawPlayerZone draws the icon + score of a single player inside its zone.
func (sv *ScoreView) drawPlayerZone(
	screen *ebiten.Image,
	p *game.Player,
	x, y, zoneWidth, zoneHeight float64,
) {
	padding := zonePaddingPx

	// Icon size is derived from zone height minus padding and reserved space.
	iconSize := zoneHeight - (two * padding) - zoneIconBottomPadding
	if iconSize < zoneMinIconSizePx {
		iconSize = zoneHeight - (two * padding)
	}

	// Draw name.
	if p.Name != "" {
		nameOpts := &text.DrawOptions{}
		nameOpts.PrimaryAlign = text.AlignCenter
		nameOpts.SecondaryAlign = text.AlignCenter
		nameOpts.ColorScale.ScaleWithColor(p.Color)
		nameOpts.GeoM.Translate(x+zoneWidth*half, y+padding*zoneNamePaddingRatio+zoneNameYOffsetPx)
		text.Draw(screen, p.Name, assets.NormalFont, nameOpts)
	}

	// Draw symbol.
	if p.Symbol != nil && p.Symbol.Image != nil {
		op := &ebiten.DrawImageOptions{}

		w, h := p.Symbol.Image.Bounds().Dx(), p.Symbol.Image.Bounds().Dy()
		scale := iconSize / float64(h)
		if w > h {
			scale = iconSize / float64(w)
		}
		op.GeoM.Scale(scale, scale)

		// Center icon horizontally inside zone.
		iconX := x + (zoneWidth-iconSize)*half
		iconY := y + padding + zoneIconExtraTopShiftPx
		op.GeoM.Translate(iconX, iconY)

		// Apply player color tint.
		op.ColorScale.ScaleWithColor(p.Color)

		// Dim non-active players.
		if sv.gameRef.Current != p && sv.gameRef.State == game.PLAYING {
			op.ColorScale.ScaleAlpha(nonActiveAlphaScale)
		}

		screen.DrawImage(p.Symbol.Image, op)
	}

	// Draw score text.
	msg := fmt.Sprintf("%d", p.Points)

	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter
	opts.ColorScale.ScaleWithColor(sv.Style.TextColor)

	// Score left-aligned inside zone.
	textX := x + padding
	textY := y + zoneHeight - padding
	opts.GeoM.Translate(textX, textY)

	text.Draw(screen, msg, assets.NormalFont, opts)
}
