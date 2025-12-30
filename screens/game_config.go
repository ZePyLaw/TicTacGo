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
	BoardSize int
	ToWin     int
	Players   []PlayerConfig
}

// DefaultGameConfig returns a ready-to-play configuration with 2 humans.
func DefaultGameConfig() GameConfig {
	return GameConfig{
		BoardSize: 3,
		ToWin:     3,
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
			{
				Name:   "Player 3",
				Color:  color.RGBA{R: 75, G: 192, B: 192, A: 255},
				Symbol: assets.TriangleSymbol,
				Ready:  false,
			},
			{
				Name:   "Player 4",
				Color:  color.RGBA{R: 255, G: 206, B: 86, A: 255},
				Symbol: assets.SquareSymbol,
				Ready:  false,
			},
		},
	}
}
