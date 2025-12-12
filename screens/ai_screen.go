package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"

	"github.com/hajimehoshi/ebiten/v2"
)

type AIScreen struct {
	host ScreenHost
	cfg  GameConfig
	btns []*ui.Button
}

func NewAIScreen(h ScreenHost, cfg GameConfig) *AIScreen {
	s := &AIScreen{
		host: h,
		cfg:  cfg,
	}

	s.btns = []*ui.Button{

		ui.NewButton("Easy", 0, -60, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.DefaultWidgetStyle,
			func() {
				cfg.AIModel = ai_models.RandomAI{}
				h.SetScreen(NewGameScreen(h, cfg))
			},
		),

		ui.NewButton("Hard", 0, 0, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.DefaultWidgetStyle,
			func() {
				cfg.AIModel = ai_models.MinimaxAI{}
				h.SetScreen(NewGameScreen(h, cfg))
			},
		),

		ui.NewButton("Back", 0, 80, uiutils.AnchorCenter,
			buttonWidth, buttonHeight, buttonRadius,
			uiutils.TransparentWidgetStyle,
			func() {
				h.SetScreen(NewStartScreen(h))
			},
		),
	}

	return s
}

func (s *AIScreen) Update() error {
	for _, b := range s.btns {
		b.Update()
	}
	return nil
}

func (s *AIScreen) Draw(screen *ebiten.Image) {
	for _, b := range s.btns {
		b.Draw(screen)
	}
}
