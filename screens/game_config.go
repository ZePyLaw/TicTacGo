package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"image/color"
)

// PlayerConfig contains the customization options for one player slot.
type PlayerConfig struct {
	Name    string
	Color   color.Color
	Symbol  assets.SymbolType
	IsAI    bool
	AIModel ai_models.AIModel
	Ready   bool
}

// GameConfig aggregates the custom setup before launching a match.
type GameConfig struct {
	BoardWidth  int // Number of columns in the grid
	BoardHeight int // Number of rows in the grid
	ToWin       int // Number of aligned symbols required to win
	Players     []PlayerConfig
}

// DefaultGameConfig returns a ready-to-play configuration with 2 humans.
func DefaultGameConfig() GameConfig {
	return GameConfig{
		BoardWidth:  3,
		BoardHeight: 3,
		ToWin:       3,
		Players: []PlayerConfig{
			{
				Name:   "Player 1",
				Color:  color.RGBA{R: 255, G: 99, B: 132, A: 255},
				Symbol: assets.CircleSymbol,
				Ready:  false,
			},
			{
				Name:   "Player 2",
				Color:  color.RGBA{R: 54, G: 162, B: 235, A: 255},
				Symbol: assets.CrossSymbol,
				Ready:  false,
			},
		},
	}
}
