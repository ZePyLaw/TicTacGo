package game

// Move represents a single move on the game board.
//
// X and Y are zero-based coordinates referring to a cell in the board
// grid (accessed as Cells[X][Y]).
type Move struct {
	X int // Column index
	Y int // Row index
}
