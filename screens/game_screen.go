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
	playerAI  map[*game.Player]ai_models.AIModel
}

const (
	// Minimum allowed board dimension. The UI/gameplay expects at least 3x3.
	minBoardDimension = 3

	// Board visual size in pixels.
	boardPixelSize = 480.0

	// Score view size in pixels.
	scorePixelWidth  = 300
	scorePixelHeight = 80

	// Sentinel for "no move".
	noMoveCoord = -1

	// Number of frames a key must be held to trigger global action.
	keyHoldFramesToTrigger = 60

	// Opaque alpha channel value.
	colorAlphaOpaque = 255
)

var (
	// Color used for the end-of-game message (yellow).
	endMessageColor = color.RGBA{R: 255, G: 255, B: 0, A: colorAlphaOpaque}

	// Default colors used when player config does not define a color.
	defaultPlayerColors = []color.Color{
		color.RGBA{R: 255, G: 99, B: 132, A: colorAlphaOpaque},
		color.RGBA{R: 54, G: 162, B: 235, A: colorAlphaOpaque},
		color.RGBA{R: 75, G: 192, B: 192, A: colorAlphaOpaque},
		color.RGBA{R: 255, G: 206, B: 86, A: colorAlphaOpaque},
	}
)

// NewGameScreen initializes a new GameScreen with a fresh game and board view.
func NewGameScreen(h ScreenHost, cfg GameConfig) *GameScreen {
	boardWidth := cfg.BoardWidth
	if boardWidth < minBoardDimension {
		boardWidth = minBoardDimension
	}
	boardHeight := cfg.BoardHeight
	if boardHeight < minBoardDimension {
		boardHeight = minBoardDimension
	}

	// ToWin must be <= smallest dimension
	minDim := boardWidth
	if boardHeight < minDim {
		minDim = boardHeight
	}
	toWin := cfg.ToWin
	if toWin <= 0 || toWin > minDim {
		toWin = minDim
	}

	players, aiMap := buildPlayers(cfg)

	// Create game logic
	g := game.NewGameWithConfig(boardWidth, boardHeight, toWin, players)

	gs := &GameScreen{
		host:     h,
		game:     g,
		playerAI: aiMap,
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
			model := gs.playerAI[current]
			if model != nil {
				x, y := model.NextMove(gs.game.Board, current, gs.game.Players)
				if x >= 0 && y >= 0 && x != noMoveCoord && y != noMoveCoord {
					gs.game.PlayMove(x, y)
				}
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
	if inpututil.KeyPressDuration(ebiten.KeyEscape) == keyHoldFramesToTrigger {
		os.Exit(0)
	}
	if inpututil.KeyPressDuration(ebiten.KeyR) == keyHoldFramesToTrigger {
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
		msg = fmt.Sprintf("%s wins!", gs.game.Winner.Name)
	} else {
		msg = "It's a draw!"
	}

	opts := &text.DrawOptions{}
	opts.PrimaryAlign = text.AlignCenter
	opts.SecondaryAlign = text.AlignCenter

	sw, sh := screen.Bounds().Dx(), screen.Bounds().Dy()
	opts.GeoM.Translate(float64(sw)/2, float64(sh)/2)

	opts.ColorScale.ScaleWithColor(endMessageColor)
	text.Draw(screen, msg, assets.BigFont, opts)
}

// buildPlayers turns the setup configuration into runtime players
// and returns a map of AI models keyed by player for quick lookup.
func buildPlayers(cfg GameConfig) ([]*game.Player, map[*game.Player]ai_models.AIModel) {
	var players []*game.Player
	aiByPlayer := map[*game.Player]ai_models.AIModel{}

	readyCount := 0
	for _, pc := range cfg.Players {
		if pc.Ready {
			readyCount++
		}
	}

	colorIdx := 0
	for idx, pc := range cfg.Players {
		if readyCount > 0 && !pc.Ready {
			continue
		}

		c := pc.Color
		if c == nil {
			c = defaultPlayerColors[colorIdx%len(defaultPlayerColors)]
		}
		colorIdx++

		sym := assets.NewSymbol(pc.Symbol)
		p := game.NewPlayer(sym, c)
		if pc.Name != "" {
			p.Name = pc.Name
		} else {
			p.Name = fmt.Sprintf("Player %d", idx+1)
		}
		p.IsAI = pc.IsAI

		players = append(players, p)

		if pc.IsAI {
			model := pc.AIModel
			if model == nil {
				model = ai_models.RandomAI{}
			}
			aiByPlayer[p] = model
		}
	}

	if len(players) == 0 {
		players = []*game.Player{
			game.NewPlayer(assets.NewSymbol(assets.CircleSymbol), defaultPlayerColors[0]),
			game.NewPlayer(assets.NewSymbol(assets.CrossSymbol), defaultPlayerColors[1]),
		}
		for i, p := range players {
			p.Name = fmt.Sprintf("Player %d", i+1)
		}
	}

	return players, aiByPlayer
}
