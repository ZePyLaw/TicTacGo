package ai_models

import "GoTicTacToe/game"

type MinimaxAI struct{}

func (MinimaxAI) NextMove(board *game.Board, me *game.Player, players [2]*game.Player) (int, int) {
	bestScore := -9999
	bestMove := game.Move{X: -1, Y: -1}

	for _, mv := range board.AvailableMoves() {
		clone := board.Clone()
		clone.Play(me, mv.X, mv.Y)

		score := minimax(clone, me, players, false)

		if score > bestScore {
			bestScore = score
			bestMove = mv
		}
	}

	return bestMove.X, bestMove.Y
}

func minimax(board *game.Board, me *game.Player, players [2]*game.Player, maximizing bool) int {
	winner := board.CheckWin()
	if winner == me {
		return +1
	}
	if winner != nil && winner != me {
		return -1
	}
	if board.CheckDraw() {
		return 0
	}

	if maximizing {
		best := -9999
		for _, mv := range board.AvailableMoves() {
			clone := board.Clone()
			clone.Play(me, mv.X, mv.Y)
			score := minimax(clone, me, players, false)
			if score > best {
				best = score
			}
		}
		return best
	}

	// minimizing (opponent turn)
	opp := me.Opponent(players)
	best := 9999

	for _, mv := range board.AvailableMoves() {
		clone := board.Clone()
		clone.Play(opp, mv.X, mv.Y)
		score := minimax(clone, me, players, true)
		if score < best {
			best = score
		}
	}

	return best
}
