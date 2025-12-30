package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type StartScreen struct {
	host       ScreenHost
	buttons    []*ui.Button
	buttonPane *ui.Container
	background *ebiten.Image
}

const (
	buttonWidth      = float64(200)
	buttonHeight     = float64(64)
	buttonRadius     = float64(10)
	buttonSpacing    = float64(20)
	buttonYOffset    = float64(140)
	buttonPaneWidth  = float64(720)
	buttonPaneHeight = float64(200)
)

func NewStartScreen(h ScreenHost) *StartScreen {
	s := &StartScreen{host: h}

	s.buttonPane = ui.NewContainer(
		0, buttonYOffset,
		buttonPaneWidth, buttonPaneHeight,
		uiutils.AnchorCenter,
		12,
		uiutils.TransparentWidgetStyle,
	)
	s.buttonPane.Padding = uiutils.InsetsAll(12)

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

		ui.NewButton(
			"Customize",
			buttonWidth+buttonSpacing, 0,
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

func (s *StartScreen) Update() error {
	s.buttonPane.Update()
	return nil
}

func (s *StartScreen) Draw(screen *ebiten.Image) {
	w, h := screen.Bounds().Dx(), screen.Bounds().Dy()
	if s.background == nil {
		topColor := color.RGBA{R: 0x1C, G: 0x25, B: 0x41, A: 0xFF}
		bottomColor := color.RGBA{R: 0x0B, G: 0x13, B: 0x2B, A: 0xFF}
		s.background = uiutils.CreateGradientBackground(w, h, topColor, bottomColor)
	}
	screen.DrawImage(s.background, nil)

	logo := assets.Logo
	origLogoW := float64(logo.Bounds().Dx())
	origLogoH := float64(logo.Bounds().Dy())

	op := &ebiten.DrawImageOptions{}

	op.Filter = ebiten.FilterLinear

	targetWidth := float64(w) * 0.5

	scaleFactor := targetWidth / origLogoW

	op.GeoM.Scale(scaleFactor, scaleFactor)

	scaledLogoW := origLogoW * scaleFactor
	scaledLogoH := origLogoH * scaleFactor

	posX := (float64(w) - scaledLogoW) / 2
	posY := (float64(h) - scaledLogoH) / 2 * 0.4

	op.GeoM.Translate(posX, posY)

	screen.DrawImage(logo, op)

	// Buttons
	s.buttonPane.Draw(screen)
}
