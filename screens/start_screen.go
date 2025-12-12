package screens

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type StartScreen struct {
	host       ScreenHost
	buttons    []*ui.Button
	background *ebiten.Image
}

const (
	buttonWidth  = float64(200)
	buttonHeight = float64(64)
	buttonRadius = float64(10)
)

func NewStartScreen(h ScreenHost) *StartScreen {
	s := &StartScreen{host: h}

	s.buttons = []*ui.Button{
		ui.NewButton("Local Play", -buttonWidth/2-10, 200.0, uiutils.AnchorCenter, buttonWidth, buttonHeight, buttonRadius, uiutils.DefaultWidgetStyle, func() {
			s.host.SetScreen(NewGameScreen(s.host))
		}),
		ui.NewButton("Multiplayer", buttonWidth/2+10, 200.0, uiutils.AnchorCenter, buttonWidth, buttonHeight, buttonRadius, uiutils.TransparentWidgetStyle, func() {}),
	}

		ui.NewButton("Local (2 Players)", 0, -80, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius, uiutils.DefaultWidgetStyle,
			func() {
				s.host.SetScreen(NewGameScreen(s.host, GameConfig{
					Mode: LocalVsLocal,
				}))
			},
		),

		ui.NewButton("Play vs AI", 0, 0, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.DefaultWidgetStyle,
			func() {
				h.SetScreen(NewAIScreen(h, GameConfig{
					Mode: LocalVsAI,
				}))
			},
		),

		ui.NewButton(
			"Multiplayer",
			0, 80,
			uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.TransparentWidgetStyle,
			func() {},
		),
	}
	return s
}

func (s *StartScreen) Update() error {
	for _, btn := range s.buttons {
		btn.Update()
	}
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

	targetWidth := float64(w) * 0.7

	scaleFactor := targetWidth / origLogoW

	op.GeoM.Scale(scaleFactor, scaleFactor)

	scaledLogoW := origLogoW * scaleFactor
	scaledLogoH := origLogoH * scaleFactor

	posX := (float64(w) - scaledLogoW) / 2
	posY := (float64(h) - scaledLogoH) / 2 * 0.6

	op.GeoM.Translate(posX, posY)

	screen.DrawImage(logo, op)

	// Buttons
	for _, btn := range s.buttons {
		btn.Draw(screen)
	}
}
