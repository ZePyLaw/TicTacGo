package game

import (
	"GoTicTacToe/assets"
	"image/color"
)

// GameState represents the current state of a match.
type GameState int

const (
	// PLAYING indicates that a match is currently in progress.
	PLAYING GameState = iota

	// GAME_END indicates that the match has ended (win or draw).
	GAME_END
)

// Default game configuration for a classic Tic-Tac-Toe match.
const (
	DefaultBoardWidth  = 3
	DefaultBoardHeight = 3
	DefaultToWin       = 3
)

// Default player styling (symbols and colors).
const (
	defaultColorAlpha = 255
)

var (
	defaultPlayer1Color = color.RGBA{R: 255, G: 0, B: 0, A: defaultColorAlpha} // red
	defaultPlayer2Color = color.RGBA{R: 0, G: 0, B: 255, A: defaultColorAlpha} // blue
	defaultPlayer1Sym   = assets.CircleSymbol
	defaultPlayer2Sym   = assets.CrossSymbol
)

// Game contains all data and logic required to run a match.
//
// It orchestrates the board, players, turn order, match state, and scoring.
// The board size and "toWin" are kept so Reset() can rebuild the board if needed.
type Game struct {
	State   GameState // Current game state (playing or ended)
	Board   *Board    // Game board instance
	Players []*Player // All players involved in the match
	Current *Player   // Player whose turn it currently is
	Winner  *Player   // Winner of the match (nil in case of draw)

	boardWidth  int
	boardHeight int
	toWin       int
}

// NewGame creates a new Game instance using the default configuration
// (classic 3x3 Tic-Tac-Toe, 3 aligned symbols to win) with default players.
func NewGame() *Game {
	return NewGameWithConfig(DefaultBoardWidth, DefaultBoardHeight, DefaultToWin, nil)
}

// NewGameWithConfig creates a new Game instance with a custom board configuration.
//
// If players is nil or empty, default players are created (2-player game).
func NewGameWithConfig(boardWidth, boardHeight, toWin int, players []*Player) *Game {
	g := &Game{}
	g.ResetHardWithPlayers(boardWidth, boardHeight, toWin, players)
	return g
}

// ResetHard fully resets the match (board, players, scores, state) using default configuration.
// This does NOT preserve any previous scores.
func (g *Game) ResetHard() {
	g.ResetHardWithPlayers(DefaultBoardWidth, DefaultBoardHeight, DefaultToWin, nil)
}

// ResetHardWithPlayers fully resets the match and applies a new configuration.
//
// This creates a new board and replaces the player list.
// Player scores are reset to zero.
func (g *Game) ResetHardWithPlayers(boardWidth, boardHeight, toWin int, players []*Player) {
	g.boardWidth = boardWidth
	g.boardHeight = boardHeight
	g.toWin = toWin
	g.Board = NewBoard(boardWidth, boardHeight, toWin)

	// Fallback to default players if none provided.
	if len(players) == 0 {
		players = []*Player{
			NewPlayer(assets.NewSymbol(defaultPlayer1Sym), defaultPlayer1Color),
			NewPlayer(assets.NewSymbol(defaultPlayer2Sym), defaultPlayer2Color),
		}
	}

	g.Players = players

	// Initialize players' scores to zero.
	for _, p := range g.Players {
		p.Points = 0
	}

	g.Current = g.Players[0]
	g.Winner = nil
	g.State = PLAYING
}

// Reset clears the board and restarts the match while keeping player scores intact.
//
// This is typically used between rounds in the same session.
func (g *Game) Reset() {
	if g.Board == nil {
		g.Board = NewBoard(g.boardWidth, g.boardHeight, g.toWin)
	} else {
		g.Board.Clear()
	}

	// Reset match state.
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
//
// If the current player is the last in the list, it wraps around.
// If Current is not found (unexpected state), it falls back to the first player.
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

	// Fallback: if current player is not found, reset to first player.
	g.Current = g.Players[0]
}

// PlayMove attempts to play a move at (x, y), then updates the match state.
//
// It handles:
// - move validation (via Board.Play)
// - win detection and scoring
// - draw detection
// - switching to the next player when the match continues
func (g *Game) PlayMove(x, y int) bool {
	ok := g.Board.Play(g.Current, x, y)
	if !ok {
		return false
	}

	// Check for victory.
	if g.CheckWin() {
		return true
	}

	// Check for draw (board full, no winner).
	if g.CheckDraw() {
		return true
	}

	// Continue the game: switch to next player.
	g.NextPlayer()
	return true
}

// CheckWin checks whether a player won the match.
//
// If a winner is found, it updates Winner, increments the winner's score,
// and ends the match (State = GAME_END).
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

// CheckDraw checks whether the match is a draw (board full, no winner).
//
// If the match is a draw, Winner is set to nil and the match ends
// (State = GAME_END).
func (g *Game) CheckDraw() bool {
	if g.Board.CheckDraw() {
		g.Winner = nil
		g.State = GAME_END
		return true
	}
	return false
}
