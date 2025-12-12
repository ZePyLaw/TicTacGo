package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"GoTicTacToe/game"
	"GoTicTacToe/ui"
	uiutils "GoTicTacToe/ui/utils"
	"fmt"
	"image/color"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// GameScreen represents the main gameplay view.
// It handles the game logic, UI board rendering, and HUD display.
type GameScreen struct {
	host      ScreenHost
	game      *game.Game
	boardView *ui.BoardView
	scoreView *ui.ScoreView
	players   [2]*game.Player
	aimodels  [2]ai_models.AIModel
}

const (
	boardPixelSize   = 480.0 // Board visual size in pixels
	scorePixelWidth  = 300
	scorePixelHeight = 80
)

// NewGameScreen initializes a new GameScreen with a fresh game and board view.
func NewGameScreen(h ScreenHost, cfg GameConfig) *GameScreen {

	// Create game logic
	g := game.NewGame()

	gs := &GameScreen{
		host: h,
		game: g,
	}

	// Store player pointers
	gs.players[0] = gs.game.Players[0]
	gs.players[1] = gs.game.Players[1]

	// Apply configuration
	if cfg.Mode == LocalVsAI && cfg.AIModel != nil {
		gs.players[1].IsAI = true
		gs.aimodels[1] = cfg.AIModel
	}

	gs.scoreView = ui.NewScoreView(g, scorePixelWidth, scorePixelHeight, uiutils.DefaultWidgetStyle)

	// Create the interactive board view with callback on cell click
	gs.boardView = ui.NewBoardView(
		g.Board, // Logical board reference
		0, 0,
		boardPixelSize, // Pixel size
		uiutils.DefaultWidgetStyle,
		func(x, y int) {
			gs.game.PlayMove(x, y)
		},
	)

	return gs
}

// Update processes input and updates UI components.
func (gs *GameScreen) Update() error {

	// Handle AI board interactions
	if gs.game.State == game.PLAYING {
		current := gs.game.Current
		if current.IsAI {

			// Find index of current player
			var idx int
			if gs.players[0] == current {
				idx = 0
			} else {
				idx = 1
			}

			model := gs.aimodels[idx]
			if model != nil {
				x, y := model.NextMove(gs.game.Board, current, gs.players)
				gs.game.PlayMove(x, y)
			}

			return nil // skip human input this frame
		}
	}

	// Handle Human board interactions
	gs.boardView.Update()

	// Reset the game if it's finished and the user clicks anywhere
	if gs.game.State == game.GAME_END {
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			gs.game.Reset()
		}
	}

	// Global hotkeys
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == 60 {
		os.Exit(0)
	}
	if inpututil.KeyPressDuration(ebiten.KeyR) == 60 {
		gs.game.Reset()
		gs.game.ResetPoints()
	}

	return nil
}

// Draw renders the board and HUD.
func (gs *GameScreen) Draw(screen *ebiten.Image) {
	// Draw board component
	gs.boardView.Draw(screen)
	gs.scoreView.Draw(screen)

	// Display win/draw message if needed
	if gs.game.State == game.GAME_END {
		gs.drawEndMessage(screen)
	}
}

// drawEndMessage displays a centered win/draw message at the end of a game.
func (gs *GameScreen) drawEndMessage(screen *ebiten.Image) {
	var msg string
	if gs.game.Winner != nil {
		msg = fmt.Sprintf("%s wins!", gs.game.Winner.Symbol.Type)
	} else {
		msg = "It's a draw!"
	}

	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	opts.GeoM.Translate(float64(sw)/2, float64(sh)/2)

	opts.ColorScale.ScaleWithColor(color.RGBA{255, 255, 0, 255})

	text.Draw(screen, msg, assets.BigFont, opts)
}
