package ui

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// cardSymbolCache stores pre-rendered symbol images to avoid regenerating them.
var cardSymbolCache = map[assets.SymbolType]*ebiten.Image{}

// PlayerCardConfig holds the configuration data for updating a player card.
// This struct is used to pass player state from the setup screen to the card.
type PlayerCardConfig struct {
	Name     string            // Display name of the player
	Subtitle string            // Secondary text (e.g., "Human" or "AI (Easy)")
	Symbol   assets.SymbolType // The symbol type this player uses
	Color    color.Color       // The player's display color
	Ready    bool              // Whether the player is ready to start
}

// PlayerCardView is a widget that displays player information in a card format.
// It shows the player's name, role, symbol, and color, with interactive elements
// for symbol selection and ready state.
type PlayerCardView struct {
	Widget // Embedded base widget providing position, size, and anchor

	Title           string            // Primary text displayed at the top of the card
	Subtitle        string            // Secondary text displayed below the title
	CenterLabel     string            // Text displayed in the center (when ShowCenterLabel is true)
	Color           color.Color       // Accent color for the card (used for strip and symbol)
	Symbol          assets.SymbolType // The symbol to display in the card center
	ShowCenterLabel bool              // If true, shows CenterLabel instead of symbol

	OnSymbolClick func() // Callback invoked when the symbol area is clicked

	ReadyButton *Button // Reference to the associated ready button for style updates

	accentStrip *ebiten.Image // Pre-rendered colored strip at the top of the card
}

// NewPlayerCard creates a new player card widget at the specified position and size.
// The card includes a rounded background and an accent strip that will be colored
// based on the player's color.
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
	// Create the rounded rectangle background
	card.image = utils.CreateRoundedRect(int(width), int(height), 14, color.RGBA{R: 24, G: 34, B: 58, A: 230})
	// Create the accent strip (will be filled with player color during draw)
	card.accentStrip = ebiten.NewImage(int(width), 14)
	return card
}

// Update handles input events for the player card.
// Currently only handles symbol click detection when the symbol is visible.
func (c *PlayerCardView) Update() {
	// Skip if showing center label or no click handler
	if c.ShowCenterLabel || c.OnSymbolClick == nil {
		return
	}

	// Check for mouse click release
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		rect := c.LayoutRect()
		iconX, iconY, iconSize := c.symbolRect(rect)
		mx, my := ebiten.CursorPosition()

		// Check if click is within symbol bounds
		if float64(mx) >= iconX && float64(mx) <= iconX+iconSize &&
			float64(my) >= iconY && float64(my) <= iconY+iconSize {
			c.OnSymbolClick()
		}
	}
}

// Draw renders the player card to the screen.
// The card consists of: background, accent strip, title/subtitle text, and symbol.
func (c *PlayerCardView) Draw(screen *ebiten.Image) {
	rect := c.LayoutRect()

	// Draw the card background
	if c.image != nil {
		op := &ebiten.DrawImageOptions{}
		scaleX := rect.Width / float64(c.image.Bounds().Dx())
		scaleY := rect.Height / float64(c.image.Bounds().Dy())
		op.GeoM.Scale(scaleX, scaleY)
		op.GeoM.Translate(rect.X, rect.Y)
		screen.DrawImage(c.image, op)
	}

	// If showing center label, draw only that and return
	if c.ShowCenterLabel {
		c.drawCenterLabel(screen, rect)
		return
	}

	// Draw the colored accent strip at the top
	if c.accentStrip != nil {
		accent := c.colorOr(color.White)
		accentRGBA := color.RGBAModel.Convert(accent).(color.RGBA)
		accentRGBA.A = 210
		c.accentStrip.Fill(accentRGBA)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(rect.X, rect.Y)
		screen.DrawImage(c.accentStrip, op)
	}

	// Draw title and subtitle text
	c.drawText(screen, rect)
	// Draw the player's symbol
	c.drawSymbol(screen, rect)
}

// drawText renders the title and subtitle text on the card.
func (c *PlayerCardView) drawText(screen *ebiten.Image, rect utils.LayoutRect) {
	paddingX := rect.Width * 0.05
	titleY := rect.Y + rect.Height*0.14
	subtitleY := rect.Y + rect.Height*0.29

	// Draw title (player name)
	if c.Title != "" {
		titleOpts := &text.DrawOptions{}
		titleOpts.PrimaryAlign = text.AlignStart
		titleOpts.SecondaryAlign = text.AlignCenter
		titleOpts.ColorScale.ScaleWithColor(color.White)
		titleOpts.GeoM.Translate(rect.X+paddingX, titleY)
		text.Draw(screen, c.Title, assets.NormalFont, titleOpts)
	}

	// Draw subtitle (role description)
	if c.Subtitle != "" {
		subOpts := &text.DrawOptions{}
		subOpts.PrimaryAlign = text.AlignStart
		subOpts.SecondaryAlign = text.AlignCenter
		subOpts.ColorScale.ScaleWithColor(color.RGBA{R: 180, G: 200, B: 230, A: 255})
		subOpts.GeoM.Translate(rect.X+paddingX, subtitleY)
		text.Draw(screen, c.Subtitle, assets.NormalFont, subOpts)
	}
}

// drawSymbol renders the player's symbol in the center of the card.
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

	// Scale to fit within the icon size while maintaining aspect ratio
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

// drawCenterLabel renders centered text in the middle of the card.
// Used as an alternative to showing the symbol.
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

// symbolRect calculates the position and size of the symbol area within the card.
// Returns (x, y, size) where size is used for both width and height.
func (c *PlayerCardView) symbolRect(rect utils.LayoutRect) (float64, float64, float64) {
	iconSize := rect.Height * 0.35
	iconX := rect.X + (rect.Width-iconSize)/2
	iconY := rect.Y + rect.Height*0.4
	return iconX, iconY, iconSize
}

// colorOr returns the card's color if set, otherwise returns the fallback color.
func (c *PlayerCardView) colorOr(fallback color.Color) color.Color {
	if c.Color != nil {
		return c.Color
	}
	return fallback
}

// cachedCardSymbol retrieves or creates a cached image for the given symbol type.
// This avoids regenerating symbol images on every frame.
func cachedCardSymbol(sym assets.SymbolType) *ebiten.Image {
	if img, ok := cardSymbolCache[sym]; ok && img != nil {
		return img
	}
	img := assets.NewSymbol(sym).Image
	cardSymbolCache[sym] = img
	return img
}

// UpdateFromConfig updates the card's display properties from a PlayerCardConfig.
// This method also updates the associated ReadyButton's label and style.
func (c *PlayerCardView) UpdateFromConfig(cfg PlayerCardConfig) {
	c.Title = cfg.Name
	c.Subtitle = cfg.Subtitle
	c.Symbol = cfg.Symbol
	c.Color = cfg.Color
	c.ShowCenterLabel = false

	// Update the ready button if one is associated
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
