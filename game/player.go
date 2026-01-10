package game

import (
	"GoTicTacToe/assets"
	"image/color"
)

// Player represents a participant in the game.
//
// A player can be either human-controlled or AI-controlled.
// The Symbol and Color are used for rendering, while Points tracks
// the player's score across multiple rounds.
type Player struct {
	Symbol *assets.Symbol // Visual symbol associated with the player
	Points int            // Score accumulated across rounds
	Color  color.Color    // Display color used in the UI
	Name   string         // Optional player name
	IsAI   bool           // Indicates whether the player is AI-controlled
}

// NewPlayer creates and returns a new player instance.
//
// The player's score is initialized to zero.
// The Name and IsAI fields may be set later by the caller if needed.
func NewPlayer(sym *assets.Symbol, color color.Color) *Player {
	return &Player{
		Symbol: sym,
		Points: 0,
		Color:  color,
	}
}

// Opponent returns the opponent of the current player.
//
// This method assumes a two-player game and returns the first player
// in the list that is not the receiver. If no opponent is found,
// nil is returned.
//
// Note: For games with more than two players, this logic would need
// to be adapted.
func (p *Player) Opponent(players []*Player) *Player {
	for _, candidate := range players {
		if candidate != p {
			return candidate
		}
	}
	return nil
}
