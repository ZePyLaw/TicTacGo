package game

// Board represents the game grid and contains player tokens.
//   - Cells is a Width x Height matrix of *Player (accessed as Cells[x][y]).
//   - Width is the number of columns.
//   - Height is the number of rows.
//   - ToWin defines how many aligned symbols are required to win
//     (used for variants such as 4-in-a-row or 5-in-a-row).
type Board struct {
	Cells  [][]*Player
	Width  int // Number of columns
	Height int // Number of rows
	ToWin  int
}

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
func (b *Board) CheckWin() *Player {
	target := b.ToWin
	minDim := b.Width
	if b.Height < minDim {
		minDim = b.Height
	}
	if target <= 0 || target > minDim {
		target = minDim
	}

	directions := [][2]int{
		{1, 0},  // horizontal
		{0, 1},  // vertical
		{1, 1},  // diagonal down-right
		{1, -1}, // diagonal up-right
	}

	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			start := b.Cells[x][y]
			if start == nil {
				continue
			}

			for _, dir := range directions {
				count := 1
				for step := 1; step < target; step++ {
					nx := x + dir[0]*step
					ny := y + dir[1]*step
					if nx < 0 || ny < 0 || nx >= b.Width || ny >= b.Height {
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

// CheckDraw returns true if the board is full and no winner exists.
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
	moves := []Move{}
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
func (b *Board) Clone() *Board {
	clone := NewBoard(b.Width, b.Height, b.ToWin)
	for x := 0; x < b.Width; x++ {
		for y := 0; y < b.Height; y++ {
			clone.Cells[x][y] = b.Cells[x][y]
		}
	}
	return clone
}
