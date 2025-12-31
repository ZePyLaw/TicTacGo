package ui

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/game"
	"GoTicTacToe/ui/utils"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// ScoreView displays player icons and scores for ANY number of players.
type ScoreView struct {
	Widget
	gameRef *game.Game
}

// NewScoreView creates a flexible score panel widget.
func NewScoreView(g *game.Game, width, height float64, style utils.WidgetStyle) *ScoreView {
	bg := utils.CreateRoundedRect(int(width), int(height), 10, style.BackgroundNormal)

	return &ScoreView{
		Widget: Widget{
			Width:   width,
			Height:  height,
			image:   bg,
			Anchor:  utils.AnchorTopCenter,
			OffsetY: 20,
			Style:   style,
		},
		gameRef: g,
	}
}

// Draw renders the score panel and each player's information.
func (sv *ScoreView) Draw(screen *ebiten.Image) {
	rect := sv.LayoutRect()
	x, y := rect.X, rect.Y

	// Draw background rectangle
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

	// Split the width into equal zones for each player
	zoneWidth := rect.Width / float64(playerCount)

	for i, p := range sv.gameRef.Players {
		zoneX := x + float64(i)*zoneWidth
		sv.drawPlayerZone(screen, p, zoneX, y, zoneWidth, rect.Height)
	}
}

// drawPlayerZone draws the icon + score of a single player inside its zone.
func (sv *ScoreView) drawPlayerZone(
	screen *ebiten.Image, p *game.Player,
	x, y, zoneWidth, zoneHeight float64,
) {
	padding := 12.0
	iconSize := zoneHeight - (2 * padding) - 16
	if iconSize < 16 {
		iconSize = zoneHeight - (2 * padding)
	}

	// --- Draw name ---
	if p.Name != "" {
		nameOpts := &text.DrawOptions{}
		nameOpts.PrimaryAlign = text.AlignCenter
		nameOpts.SecondaryAlign = text.AlignCenter
		nameOpts.ColorScale.ScaleWithColor(p.Color)
		nameOpts.GeoM.Translate(x+zoneWidth/2, y+padding/2+2)
		text.Draw(screen, p.Name, assets.NormalFont, nameOpts)
	}

	// --- Draw symbol ---
	if p.Symbol != nil && p.Symbol.Image != nil {
		op := &ebiten.DrawImageOptions{}

		w, h := p.Symbol.Image.Bounds().Dx(), p.Symbol.Image.Bounds().Dy()
		scale := iconSize / float64(h)
		if w > h {
			scale = iconSize / float64(w)
		}
		op.GeoM.Scale(scale, scale)

		// Center icon horizontally inside zone
		iconX := x + (zoneWidth-iconSize)/2
		iconY := y + padding + 12
		op.GeoM.Translate(iconX, iconY)

		// Apply player color tint
		op.ColorScale.ScaleWithColor(p.Color)

		// Dim non-active players
		if sv.gameRef.Current != p && sv.gameRef.State == game.PLAYING {
			op.ColorScale.ScaleAlpha(0.5)
		}

		screen.DrawImage(p.Symbol.Image, op)
	}

	// --- Draw score text ---
	msg := fmt.Sprintf("%d", p.Points)

	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter
	opts.ColorScale.ScaleWithColor(sv.Style.TextColor)

	// Score left-aligned inside zone
	textX := x + padding
	textY := y + zoneHeight - padding

	opts.GeoM.Translate(textX, textY)

	text.Draw(screen, msg, assets.NormalFont, opts)
}
