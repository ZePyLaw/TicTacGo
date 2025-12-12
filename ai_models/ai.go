package ai_models

import "GoTicTacToe/game"

type AIModel interface {
	NextMove(board *game.Board, me *game.Player, players [2]*game.Player) (x, y int)
}
