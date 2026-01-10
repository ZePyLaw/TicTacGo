/**
 ******************************************************************************
 * @file            : ai_model.go
 * @brief           : GoTicTacToe - AI interface definition
 * @author          : Alexandre Schmid <alexandre.schmid@edu.heia-fr.ch>
 * @author          : Jeremy Prin <alexandre.schmid@edu.heia-fr.ch>
 * @date            : 09. January 2026
 ******************************************************************************
 * @copyright   : Copyright (c) 2026 HEIA-FR / ISC
 *                Haute école d'ingénierie et d'architecture de Fribourg
 *                Informatique et Systèmes de Communication
 * @attention   : SPDX-License-Identifier: MIT OR Apache-2.0
 ******************************************************************************
 * @details
 * This file defines the common interface implemented by all AI models.
 * It allows the game engine to remain independent from the concrete
 * AI strategy used (random, minimax, or future extensions).
 ******************************************************************************
 */

package ai_models

import "GoTicTacToe/game"

// AIModel defines the behavior required for an artificial intelligence
// capable of playing the Tic-Tac-Toe game.
//
// Any AI strategy (random, minimax, etc.) must implement this interface
// in order to be usable by the game engine.
type AIModel interface {
	// NextMove computes and returns the next move to play.
	//
	// Parameters:
	// - board: current game board state
	// - me: the player controlled by this AI
	// - players: list of all players in the game
	//
	// Returns:
	// - x, y: coordinates of the chosen move on the board
	NextMove(board *game.Board, me *game.Player, players []*game.Player) (x, y int)
}
