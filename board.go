package main

import "fmt"

type Board struct {
	size int // size of the board
	cells [][]*Player // 2D array of pointers to players
	numSymbolToWin int // number of symbols in a row needed to win
}

// NewBoard creates a new board of the specified size, initializing all cells to nil.
// It returns a pointer to the newly created board.
// Example usage:
// board := NewBoard(3)
// _ | _ | _
// _ | _ | _
// _ | _ | _
func NewBoard(size int, numSymbolToWin int) *Board {
	board := &Board{
		size: size,
		cells: make([][]*Player, size),
		numSymbolToWin: numSymbolToWin,
	}

	for i := 0; i < size; i++ {
		board.cells[i] = make([]*Player, size)
	}

	return board
}

func (b *Board) play(player *Player, row, col int) error {
	if row < 0 || row >= b.size || col < 0 || col >= b.size {
		return fmt.Errorf("invalid move")
	}

	if b.cells[row][col] != nil {
		return fmt.Errorf("square already occupied")
	}

	b.cells[row][col] = player
	return nil
}

// CheckWin checks if there is a winning combination on the board.
// It returns the winning player if there is one, or nil otherwise.
// Example usage:
// winner := board.CheckWin()
// if winner != nil {
//     fmt.Printf("%s wins!\n", winner.Name())
// }
func (b *Board) CheckWin() *Player {
	// Check rows
	for row := 0; row < b.size; row++ {
		for col := 0; col <= b.size-b.numSymbolToWin; col++ {
			first := b.cells[row][col]
			if first == nil {
				continue
			}
			win := true
			for i := 1; i < b.numSymbolToWin; i++ {
				if b.cells[row][col+i] != first {
					win = false
					break
				}
			}
			if win {
				return first
			}
		}
	}

	// Check columns
	for col := 0; col < b.size; col++ {
		for row := 0; row <= b.size-b.numSymbolToWin; row++ {
			first := b.cells[row][col]
			if first == nil {
				continue
			}
			win := true
			for i := 1; i < b.numSymbolToWin; i++ {
				if b.cells[row+i][col] != first {
					win = false
					break
				}
			}
			if win {
				return first
			}
		}
	}

	// Check diagonals
	for row := 0; row <= b.size-b.numSymbolToWin; row++ {
		for col := 0; col <= b.size-b.numSymbolToWin; col++ {
			first := b.cells[row][col]
			if first == nil {
				continue
			}
			win := true
			for i := 1; i < b.numSymbolToWin; i++ {
				if b.cells[row+i][col+i] != first {
					win = false
					break
				}
			}
			if win {
				return first
			}
		}
	}

	// Check anti-diagonals
	for row := 0; row <= b.size-b.numSymbolToWin; row++ {
		for col := b.numSymbolToWin - 1; col < b.size; col++ {
			first := b.cells[row][col]
			if first == nil {
				continue
			}
			win := true
			for i := 1; i < b.numSymbolToWin; i++ {
				if b.cells[row+i][col-i] != first {
					win = false
					break
				}
			}
			if win {
				return first
			}
		}
	}

	return nil
}

// CheckDraw checks if the game is a draw.
// It returns true if there are no empty cells left, false otherwise.
// Example usage:
// isDraw := board.CheckDraw()
// if isDraw {
//     fmt.Println("It's a draw!")
// }
func (b *Board) CheckDraw() bool {
	for _, row := range b.cells {
		for _, square := range row {
			if square == nil {
				return false
			}
		}
	}
	return true
}
