package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

// StartScreen is the main menu screen.
//
// It displays the game logo and provides buttons to start a match quickly
// or open the customization setup.
type StartScreen struct {
	host       ScreenHost
	buttons    []*ui.Button
	buttonPane *ui.Container
	background *ebiten.Image
}

// Button and layout constants.
const (
	buttonWidth      = float64(200)
	buttonHeight     = float64(64)
	buttonRadius     = float64(10)
	buttonSpacing    = float64(20)
	buttonYOffset    = float64(140)
	buttonPaneWidth  = float64(720)
	buttonPaneHeight = float64(200)
)

// Start screen layout / styling constants.
const (
	// containerCornerRadiusPx is the corner radius of the button container.
	containerCornerRadiusPx = 12

	// containerPaddingPx is the padding applied inside the button container.
	containerPaddingPx = 12

	// half is used for centering computations (equivalent to dividing by 2).
	half = 0.5
)

// Visual constants for the start screen.
const (
	// logoTargetWidthRatio defines the portion of the screen width the logo should occupy.
	logoTargetWidthRatio = 0.5

	// logoVerticalBias controls how high the logo is placed (relative to vertical center).
	// 0.5 would be center; smaller values move it upwards.
	logoVerticalBias = 0.4
)

// Background gradient colors.
var (
	startBgTopColor    = color.RGBA{R: 0x1C, G: 0x25, B: 0x41, A: 0xFF}
	startBgBottomColor = color.RGBA{R: 0x0B, G: 0x13, B: 0x2B, A: 0xFF}
)

// NewStartScreen constructs the start screen and its navigation buttons.
func NewStartScreen(h ScreenHost) *StartScreen {
	s := &StartScreen{host: h}

	s.buttonPane = ui.NewContainer(
		0, buttonYOffset,
		buttonPaneWidth, buttonPaneHeight,
		uiutils.AnchorCenter,
		containerCornerRadiusPx,
		uiutils.TransparentWidgetStyle,
	)
	s.buttonPane.Padding = uiutils.InsetsAll(containerPaddingPx)

	s.buttons = []*ui.Button{
		ui.NewButton("Quick Local", -buttonWidth-buttonSpacing, 0, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius, uiutils.NormalWidgetStyle,
			func() {
				cfg := DefaultGameConfig()
				h.SetScreen(NewGameScreen(h, cfg))
			},
		),

		ui.NewButton("Quick vs AI", 0, 0, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.NormalWidgetStyle,
			func() {
				cfg := DefaultGameConfig()
				cfg.Players[1].IsAI = true
				cfg.Players[1].AIModel = ai_models.MinimaxAI{}
				h.SetScreen(NewGameScreen(h, cfg))
			},
		),

		ui.NewButton("Customize", buttonWidth+buttonSpacing, 0,
			uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.DefaultWidgetStyle,
			func() {
				h.SetScreen(NewSetupScreen(h, DefaultGameConfig()))
			},
		),
	}

	for _, btn := range s.buttons {
		s.buttonPane.AddChild(btn)
	}
	return s
}

// Update updates UI interactions for the start screen.
func (s *StartScreen) Update() error {
	s.buttonPane.Update()
	return nil
}

// Draw renders the background, logo, and buttons.
func (s *StartScreen) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()

	// Create and cache gradient background.
	if s.background == nil {
		s.background = uiutils.CreateGradientBackground(w, h, startBgTopColor, startBgBottomColor)
	}
	screen.DrawImage(s.background, nil)

	// Draw logo.
	logo := assets.Logo
	origLogoW := float64(logo.Bounds().Dx())
	origLogoH := float64(logo.Bounds().Dy())

	op := &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear

	targetWidth := float64(w) * logoTargetWidthRatio
	scaleFactor := targetWidth / origLogoW

	op.GeoM.Scale(scaleFactor, scaleFactor)

	scaledLogoW := origLogoW * scaleFactor
	scaledLogoH := origLogoH * scaleFactor

	posX := (float64(w) - scaledLogoW) * half
	posY := (float64(h) - scaledLogoH) * half * logoVerticalBias

	op.GeoM.Translate(posX, posY)
	screen.DrawImage(logo, op)

	// Draw buttons.
	s.buttonPane.Draw(screen)
}
