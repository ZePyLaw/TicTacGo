package game

import (
	"GoTicTacToe/assets"
	"image/color"
)

// GameState represents the current state of the game.
type GameState int

const (
	PLAYING  GameState = iota // A game is currently in progress
	GAME_END                  // The game has ended (win or draw)
)

// Game contains all data and logic required to run a match.
type Game struct {
	State   GameState // Current game state (playing or ended)
	Board   *Board    // Game board instance
	Players []*Player // All players involved in the game
	Current *Player   // Player whose turn it currently is
	Winner  *Player   // Winner of the match (nil in case of draw)
}

// NewGame creates a new Game instance with a full reset.
func NewGame() *Game {
	g := &Game{}
	g.ResetHard()
	return g
}

// ResetHard fully resets the game state, board, players, and winner.
// This does NOT preserve any previous scores.
func (g *Game) ResetHard() {
	g.Board = NewBoard(3, 3)

	g.Players = []*Player{
		NewPlayer(assets.NewSymbol(assets.CircleSymbol), color.RGBA{R: 255, G: 0, B: 0, A: 255}),
		NewPlayer(assets.NewSymbol(assets.CrossSymbol), color.RGBA{R: 0, G: 0, B: 255, A: 255}),
	}

	g.Current = g.Players[0]
	g.Winner = nil
	g.State = PLAYING
}

// Reset clears the board and restarts the game,
// while keeping player scores intact.
func (g *Game) Reset() {
	g.Board.Clear()
	g.Current = g.Players[0]
	g.Winner = nil
	g.State = PLAYING
}

// ResetPoints resets every player's score to zero.
func (g *Game) ResetPoints() {
	for _, p := range g.Players {
		p.Points = 0
	}
}

// NextPlayer switches the turn to the next player in the list.
// If the current player is the last in the list, it wraps around.
func (g *Game) NextPlayer() {
	if len(g.Players) == 0 {
		return
	}
	for i, p := range g.Players {
		if p == g.Current {
			g.Current = g.Players[(i+1)%len(g.Players)]
			return
		}
	}
	// Fallback: if current player is not found, reset to first player
	g.Current = g.Players[0]
}

// PlayMove attempts to play a move at (x,y), then handles
// turn switching, win detection, and draw detection.
func (g *Game) PlayMove(x, y int) bool {
	// Attempt to place the player's mark
	ok := g.Board.Play(g.Current, x, y)
	if !ok {
		return false // Invalid move (cell already taken or out of bounds)
	}

	// Check for victory
	if g.CheckWin() {
		return true
	}

	// Check for draw (board full, no winner)
	if g.CheckDraw() {
		return true
	}

	// Continue the game: switch to next player
	g.NextPlayer()
	return true
}

// CheckWin checks if a player won the match.
// If so, updates the winner, increments score, and ends the game.
func (g *Game) CheckWin() bool {
	p := g.Board.CheckWin()
	if p != nil {
		g.Winner = p
		g.Winner.Points++
		g.State = GAME_END
		return true
	}
	return false
}

// CheckDraw returns true if the game is a draw and ends the game state.
func (g *Game) CheckDraw() bool {
	if g.Board.CheckDraw() {
		g.Winner = nil
		g.State = GAME_END
		return true
	}
	return false
}

func (b *Board) AvailableMoves() []Move {
	moves := []Move{}
	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			if b.Cells[x][y] == nil {
				moves = append(moves, Move{X: x, Y: y})
			}
		}
	}
	return moves
}

func (b *Board) Clone() *Board {
	clone := NewBoard(b.Size, b.ToWin)
	for x := 0; x < b.Size; x++ {
		for y := 0; y < b.Size; y++ {
			clone.Cells[x][y] = b.Cells[x][y]
		}
	}
	return clone
}
