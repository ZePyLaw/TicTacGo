package ai_models

import (
	"GoTicTacToe/game"
	"math/rand"
)

type RandomAI struct{}

func (RandomAI) NextMove(board *game.Board, me *game.Player, players []*game.Player) (int, int) {
	moves := board.AvailableMoves()
	if len(moves) == 0 {
		return -1, -1
	}

	m := moves[rand.Intn(len(moves))]
	return m.X, m.Y
}
