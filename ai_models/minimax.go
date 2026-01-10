package ai_models

import "GoTicTacToe/game"

// MinimaxAI is an AI player using the Minimax algorithm.
// It is designed for two-player, deterministic, perfect-information games
// such as Tic-Tac-Toe.
//
// In the classic 3x3 Tic-Tac-Toe, this strategy is unbeatable (optimal play).
type MinimaxAI struct{}

// Minimax evaluation scores.
//
// These constants are intentionally kept small and symmetric:
// - win  => +1
// - loss => -1
// - draw =>  0
//
// Sentinel values are used as initial "worst possible" bounds when searching.
const (
	scoreWin  = 1
	scoreDraw = 0
	scoreLoss = -1

	// initialLowerBound is used to initialize the best score in maximizing turns.
	// It must be strictly lower than the minimal possible score (scoreLoss).
	initialLowerBound = -9999

	// initialUpperBound is used to initialize the best score in minimizing turns.
	// It must be strictly higher than the maximal possible score (scoreWin).
	initialUpperBound = 9999
)

// NextMove returns the best move (x, y) for the current player according to Minimax.
//
// The current implementation supports two-player games only.
// If the number of players is not exactly two, it falls back to RandomAI
// to avoid undefined behavior (e.g., "opponent" not well-defined).
func (MinimaxAI) NextMove(board *game.Board, me *game.Player, players []*game.Player) (int, int) {
	if len(players) != 2 {
		return RandomAI{}.NextMove(board, me, players)
	}

	bestScore := initialLowerBound
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

// minimax recursively evaluates the game tree from the perspective of "me".
//
// Parameters:
// - board: current board state
// - me: the player for which we are computing the best outcome
// - players: list of players (expected length: 2)
// - maximizing: true if it's "me" turn, false if it's opponent's turn
//
// Returns an integer score among {scoreWin, scoreDraw, scoreLoss}.
//
// This version does not implement alpha-beta pruning (which would speed up search),
// because Tic-Tac-Toe's state space is small. The code remains simple and readable.
func minimax(board *game.Board, me *game.Player, players []*game.Player, maximizing bool) int {
	// Terminal states: win/loss/draw
	winner := board.CheckWin()
	if winner == me {
		return scoreWin
	}
	if winner != nil && winner != me {
		return scoreLoss
	}
	if board.CheckDraw() {
		return scoreDraw
	}

	// Maximizing: it's "me" turn.
	if maximizing {
		best := initialLowerBound
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

	// Minimizing: opponent turn.
	opp := me.Opponent(players)
	best := initialUpperBound

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
