package game

// Board represents the game grid and contains player tokens.
//
// Cells is a Width x Height matrix of *Player (accessed as Cells[x][y]).
// Width is the number of columns.
// Height is the number of rows.
// ToWin defines how many aligned symbols are required to win (variant support).
type Board struct {
	Cells  [][]*Player
	Width  int // Number of columns
	Height int // Number of rows
	ToWin  int // Required aligned symbols to win
}

// Direction represents a 2D step (dx, dy) used for line scanning.
type Direction struct {
	DX int
	DY int
}

// Common scanning directions used for win detection.
//
// We only need to scan 4 directions to cover all possible alignments:
// horizontal, vertical, and the two diagonals.
var winDirections = [...]Direction{
	{DX: 1, DY: 0},  // horizontal (→)
	{DX: 0, DY: 1},  // vertical (↓)
	{DX: 1, DY: 1},  // diagonal down-right (↘)
	{DX: 1, DY: -1}, // diagonal up-right (↗)
}

const (
	// initialStreakCount is the count when considering the starting cell.
	initialStreakCount = 1

	// firstStep is the first step away from the starting cell.
	firstStep = 1
)

// NewBoard allocates a new empty board with given dimensions.
// All cells start as nil (empty).
func NewBoard(width, height, toWin int) *Board {
	b := &Board{
		Width:  width,
		Height: height,
		ToWin:  toWin,
		Cells:  make([][]*Player, width),
	}

	for x := range b.Cells {
		b.Cells[x] = make([]*Player, height)
	}
	return b
}

// Play attempts to place player p at grid coordinates (x, y).
// Returns true if the move is valid and the cell was empty.
func (b *Board) Play(p *Player, x, y int) bool {
	// Out-of-bounds protection
	if x < 0 || y < 0 || x >= b.Width || y >= b.Height {
		return false
	}
	// Cell already filled
	if b.Cells[x][y] != nil {
		return false
	}

	b.Cells[x][y] = p
	return true
}

// CheckWin verifies if a player has won for any streak of length ToWin
// horizontally, vertically, or diagonally.
//
// If ToWin is invalid (<= 0 or larger than the smallest board dimension),
// it is clamped to the smallest dimension. This makes the method robust
// even when used with board variants or unexpected inputs.
func (b *Board) CheckWin() *Player {
	target := b.effectiveToWin()

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			start := b.Cells[x][y]
			if start == nil {
				continue
			}

			for _, dir := range winDirections {
				count := initialStreakCount

				for step := firstStep; step < target; step++ {
					nx := x + dir.DX*step
					ny := y + dir.DY*step

					if !b.inBounds(nx, ny) {
						break
					}
					if b.Cells[nx][ny] != start {
						break
					}
					count++
				}

				if count == target {
					return start
				}
			}
		}
	}

	return nil
}

// CheckDraw returns true if the board is full (no empty cell remains).
// Note: a typical game loop should call CheckWin first; this method does not
// attempt to infer a winner.
func (b *Board) CheckDraw() bool {
	for x := range b.Cells {
		for y := range b.Cells[x] {
			if b.Cells[x][y] == nil {
				return false
			}
		}
	}
	return true
}

// Clear resets all cells to nil (empty board).
func (b *Board) Clear() {
	for x := range b.Cells {
		for y := range b.Cells[x] {
			b.Cells[x][y] = nil
		}
	}
}

// AvailableMoves returns all empty cell positions on the board.
func (b *Board) AvailableMoves() []Move {
	moves := make([]Move, 0)
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			if b.Cells[x][y] == nil {
				moves = append(moves, Move{X: x, Y: y})
			}
		}
	}
	return moves
}

// Clone creates a deep copy of the board.
//
// Note: Players are referenced (not cloned), which is intended: players are
// immutable identity objects, while the board state is what must be copied.
func (b *Board) Clone() *Board {
	clone := NewBoard(b.Width, b.Height, b.ToWin)
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			clone.Cells[x][y] = b.Cells[x][y]
		}
	}
	return clone
}

// inBounds returns true if (x, y) is within the board limits.
func (b *Board) inBounds(x, y int) bool {
	return x >= 0 && y >= 0 && x < b.Width && y < b.Height
}

// effectiveToWin returns a robust, clamped value for ToWin.
//
// If ToWin is not usable, it is clamped to min(Width, Height).
func (b *Board) effectiveToWin() int {
	target := b.ToWin
	minDim := b.Width
	if b.Height < minDim {
		minDim = b.Height
	}

	if target <= 0 || target > minDim {
		return minDim
	}
	return target
}
