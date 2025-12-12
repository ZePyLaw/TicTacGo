package screens

import "GoTicTacToe/ai_models"

type GameMode int

const (
	LocalVsLocal GameMode = iota
	LocalVsAI
)

type GameConfig struct {
	Mode    GameMode
	AIModel ai_models.AIModel // nil if not used
}
