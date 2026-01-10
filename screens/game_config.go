/**
 ******************************************************************************
 * @file            : config.go
 * @brief           : GoTicTacToe - Game and player configuration (UI setup)
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
 * This file defines the configuration structures used to customize a match
 * before it starts. These configurations are typically manipulated by the
 * user interface (screens) before creating a Game instance.
 ******************************************************************************
 */

package screens

import (
	"GoTicTacToe/ai_models"
	"GoTicTacToe/assets"
	"image/color"
)

// Default game configuration values.
const (
	defaultBoardWidth  = 3
	defaultBoardHeight = 3
	defaultToWin       = 3

	defaultColorAlpha = 255
)

// Default player colors.
var (
	defaultPlayer1Color = color.RGBA{R: 255, G: 99, B: 132, A: defaultColorAlpha}
	defaultPlayer2Color = color.RGBA{R: 54, G: 162, B: 235, A: defaultColorAlpha}
)

// PlayerConfig contains the customization options for one player slot.
//
// This structure is used by the UI to configure players before a match starts.
// It supports both human and AI-controlled players.
type PlayerConfig struct {
	Name    string            // Display name of the player
	Color   color.Color       // Player color used in the UI
	Symbol  assets.SymbolType // Symbol associated with the player
	IsAI    bool              // Indicates whether the player is AI-controlled
	AIModel ai_models.AIModel // AI strategy (used only if IsAI is true)
	Ready   bool              // Indicates whether the player is ready to start
}

// GameConfig aggregates the full setup required before launching a match.
//
// It defines the board dimensions, the win condition, and all participating
// players.
type GameConfig struct {
	BoardWidth  int            // Number of columns in the grid
	BoardHeight int            // Number of rows in the grid
	ToWin       int            // Number of aligned symbols required to win
	Players     []PlayerConfig // Player configurations
}

// DefaultGameConfig returns a ready-to-play configuration.
//
// The default configuration represents a classic 3x3 Tic-Tac-Toe game
// with two human players.
func DefaultGameConfig() GameConfig {
	return GameConfig{
		BoardWidth:  defaultBoardWidth,
		BoardHeight: defaultBoardHeight,
		ToWin:       defaultToWin,
		Players: []PlayerConfig{
			{
				Name:   "Player 1",
				Color:  defaultPlayer1Color,
				Symbol: assets.CircleSymbol,
				Ready:  false,
			},
			{
				Name:   "Player 2",
				Color:  defaultPlayer2Color,
				Symbol: assets.CrossSymbol,
				Ready:  false,
			},
		},
	}
}
