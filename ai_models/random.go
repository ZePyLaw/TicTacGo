/**
 ******************************************************************************
 * @file            : random_ai.go
 * @brief           : GoTicTacToe - Random AI implementation
 * @author          : Alexandre Schmid <alexandre.schmid@edu.heia-fr.ch>
 * @author          : Jeremy Prin <jeremy.prin@edu.heia-fr.ch>
 * @date            : 09. January 2026
 ******************************************************************************
 * @copyright   : Copyright (c) 2026 HEIA-FR / ISC
 *                Haute école d'ingénierie et d'architecture de Fribourg
 *                Informatique et Systèmes de Communication
 * @attention   : SPDX-License-Identifier: MIT OR Apache-2.0
 ******************************************************************************
 * @details
 * This file contains a simple AI that selects a random valid move.
 * It serves as:
 *   - a baseline opponent for testing,
 *   - a fallback AI when advanced strategies are not applicable,
 *   - a comparison point for more advanced algorithms (e.g. Minimax).
 ******************************************************************************
 */

package ai_models

import (
	"GoTicTacToe/game"
	"math/rand"
)

// RandomAI is an AI player that selects a move uniformly at random
// among all currently available moves.
//
// This strategy does not attempt to play optimally and can easily
// be beaten by a human player or a Minimax-based AI.
type RandomAI struct{}

// Sentinel values returned when no valid move exists.
const (
	noMoveX = -1
	noMoveY = -1
)

// NextMove returns a random valid move (x, y).
//
// If the board has no available moves left (game already finished),
// the sentinel coordinates (noMoveX, noMoveY) are returned.
func (RandomAI) NextMove(board *game.Board, _ *game.Player, _ []*game.Player) (int, int) {
	moves := board.AvailableMoves()
	if len(moves) == 0 {
		return noMoveX, noMoveY
	}

	m := moves[rand.Intn(len(moves))]
	return m.X, m.Y
}
