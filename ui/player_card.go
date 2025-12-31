package ui

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var cardSymbolCache = map[assets.SymbolType]*ebiten.Image{}

// PlayerCardConfig holds the configuration data for updating a player card.
type PlayerCardConfig struct {
	Name     string
	Subtitle string
	Symbol   assets.SymbolType
	Color    color.Color
	Ready    bool
}

// PlayerCardView is a small widget that draws a player card.
type PlayerCardView struct {
	Widget

	Title           string
	Subtitle        string
	CenterLabel     string
	Color           color.Color
	Symbol          assets.SymbolType
	ShowCenterLabel bool

	OnSymbolClick func()

	// ReadyButton is the button associated with this card for ready state.
	ReadyButton *Button

	accentStrip *ebiten.Image
}

// NewPlayerCard creates a player card widget with a fixed size.
func NewPlayerCard(
	offsetX, offsetY, width, height float64,
	anchor utils.Anchor,
) *PlayerCardView {
	card := &PlayerCardView{
		Widget: Widget{
			OffsetX: offsetX,
			OffsetY: offsetY,
			Width:   width,
			Height:  height,
			Anchor:  anchor,
		},
	}
	card.image = utils.CreateRoundedRect(int(width), int(height), 14, color.RGBA{R: 24, G: 34, B: 58, A: 230})
	card.accentStrip = ebiten.NewImage(int(width), 14)
	return card
}

// Update handles symbol click detection.
func (c *PlayerCardView) Update() {
	if c.ShowCenterLabel || c.OnSymbolClick == nil {
		return
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		rect := c.LayoutRect()
		iconX, iconY, iconSize := c.symbolRect(rect)
		mx, my := ebiten.CursorPosition()
		if float64(mx) >= iconX && float64(mx) <= iconX+iconSize &&
			float64(my) >= iconY && float64(my) <= iconY+iconSize {
			c.OnSymbolClick()
		}
	}
}

// Draw renders the card background, accent, text, and symbol.
func (c *PlayerCardView) Draw(screen *ebiten.Image) {
	rect := c.LayoutRect()

	if c.image != nil {
		op := &ebiten.DrawImageOptions{}
		scaleX := rect.Width / float64(c.image.Bounds().Dx())
		scaleY := rect.Height / float64(c.image.Bounds().Dy())
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(rect.X, rect.Y)
		screen.DrawImage(c.image, op)
	}

	if c.ShowCenterLabel {
		c.drawCenterLabel(screen, rect)
		return
	}

	if c.accentStrip != nil {
		accent := c.colorOr(color.White)
		accentRGBA := color.RGBAModel.Convert(accent).(color.RGBA)
		accentRGBA.A = 210
		c.accentStrip.Fill(accentRGBA)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(rect.X, rect.Y)
		screen.DrawImage(c.accentStrip, op)
	}

	c.drawText(screen, rect)
	c.drawSymbol(screen, rect)
}

func (c *PlayerCardView) drawText(screen *ebiten.Image, rect utils.LayoutRect) {
	paddingX := rect.Width * 0.05
	titleY := rect.Y + rect.Height*0.14
	subtitleY := rect.Y + rect.Height*0.29

	if c.Title != "" {
		titleOpts := &text.DrawOptions{}
		titleOpts.PrimaryAlign = text.AlignStart
		titleOpts.SecondaryAlign = text.AlignCenter
		titleOpts.ColorScale.ScaleWithColor(color.White)
		titleOpts.GeoM.Translate(rect.X+paddingX, titleY)
		text.Draw(screen, c.Title, assets.NormalFont, titleOpts)
	}

	if c.Subtitle != "" {
		subOpts := &text.DrawOptions{}
		subOpts.PrimaryAlign = text.AlignStart
		subOpts.SecondaryAlign = text.AlignCenter
		subOpts.ColorScale.ScaleWithColor(color.RGBA{R: 180, G: 200, B: 230, A: 255})
		subOpts.GeoM.Translate(rect.X+paddingX, subtitleY)
		text.Draw(screen, c.Subtitle, assets.NormalFont, subOpts)
	}
}

func (c *PlayerCardView) drawSymbol(screen *ebiten.Image, rect utils.LayoutRect) {
	img := cachedCardSymbol(c.Symbol)
	if img == nil {
		return
	}

	iconX, iconY, iconSize := c.symbolRect(rect)
	srcW := float64(img.Bounds().Dx())
	srcH := float64(img.Bounds().Dy())
	if srcW == 0 || srcH == 0 {
		return
	}

	scale := iconSize / srcH
	if srcW > srcH {
		scale = iconSize / srcW
	}

	symOp := &ebiten.DrawImageOptions{}
	symOp.Filter = ebiten.FilterLinear
	symOp.GeoM.Scale(scale, scale)
	symOp.GeoM.Translate(iconX, iconY)
	symOp.ColorScale.ScaleWithColor(c.colorOr(color.White))
	screen.DrawImage(img, symOp)
}

func (c *PlayerCardView) drawCenterLabel(screen *ebiten.Image, rect utils.LayoutRect) {
	label := c.CenterLabel
	if label == "" {
		return
	}

	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter
	opts.ColorScale.ScaleWithColor(color.RGBA{R: 210, G: 230, B: 255, A: 255})
	opts.GeoM.Translate(rect.X+rect.Width/2, rect.Y+rect.Height/2)
	text.Draw(screen, label, assets.NormalFont, opts)
}

func (c *PlayerCardView) symbolRect(rect utils.LayoutRect) (float64, float64, float64) {
	iconSize := rect.Height * 0.35
	iconX := rect.X + (rect.Width-iconSize)/2
	iconY := rect.Y + rect.Height*0.4
	return iconX, iconY, iconSize
}

func (c *PlayerCardView) colorOr(fallback color.Color) color.Color {
	if c.Color != nil {
		return c.Color
	}
	return fallback
}

func cachedCardSymbol(sym assets.SymbolType) *ebiten.Image {
	if img, ok := cardSymbolCache[sym]; ok && img != nil {
		return img
	}
	img := assets.NewSymbol(sym).Image
	cardSymbolCache[sym] = img
	return img
}

// UpdateFromConfig updates the card's display properties from a PlayerCardConfig.
func (c *PlayerCardView) UpdateFromConfig(cfg PlayerCardConfig) {
	c.Title = cfg.Name
	c.Subtitle = cfg.Subtitle
	c.Symbol = cfg.Symbol
	c.Color = cfg.Color
	c.ShowCenterLabel = false

	if c.ReadyButton != nil {
		if cfg.Ready {
			c.ReadyButton.Label = "Ready"
			c.ReadyButton.Style = utils.SuccessWidgetStyle
		} else {
			c.ReadyButton.Label = "Not Ready"
			c.ReadyButton.Style = utils.DefaultWidgetStyle
		}
	}
}
