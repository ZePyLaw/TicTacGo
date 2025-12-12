package screens

import (
	"GoTicTacToe/assets"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type StartScreen struct {
	host    ScreenHost
	buttons []*ui.Button
}

const (
	buttonWidth  = float64(200)
	buttonHeight = float64(64)
	buttonRadius = float64(10)
)

func NewStartScreen(h ScreenHost) *StartScreen {
	s := &StartScreen{host: h}

	s.buttons = []*ui.Button{

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
	w, h := GetWindowSize()

	// Title
	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.GeoM.Translate(w/2, h*0.2)
	opts.ColorScale.ScaleWithColor(color.White)

	text.Draw(screen, "Go Tic Tac Toe", assets.BigFont, opts)

	// Buttons
	for _, btn := range s.buttons {
		btn.Draw(screen)
	}
}
