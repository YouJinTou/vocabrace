package scrabble

import (
	"encoding/json"
	"fmt"
	"sort"
)

// BoardWidth is the standard board width.
const BoardWidth = 15

// BoardHeight is the standard board height.
const BoardHeight = 15

// BoardMinIndex is the first cell index available.
const BoardMinIndex = 0

// BoardMaxIndex is the last cell index available.
const BoardMaxIndex = 224

// BoardOrigin is the index of the center of the board.
const BoardOrigin = BoardMaxIndex / 2

// Board is a 15x15 field of cells.
type Board struct {
	Cells []Cell `json:"c"`
}

// NewBoard creates a board.
func NewBoard() *Board {
	board := Board{
		Cells: []Cell{},
	}

	return &board
}

// GetAt returns the cell at a given index or nil if not found.
func (b *Board) GetAt(i int) *Cell {
	for _, c := range b.Cells {
		if c.Index == i {
			return &c
		}
	}
	return nil
}

// SetAt sets the cell at a given index.
func (b *Board) SetAt(i int, v Cell) {
	for idx, c := range b.Cells {
		if c.Index == i {
			b.Cells[idx] = v
			return
		}
	}
}

// GetRowMinCol returns the index of the first column of a row given a cell index.
func (b *Board) GetRowMinCol(i int) int {
	row := i / BoardHeight
	minCol := row * BoardWidth
	return minCol
}

// GetRowMaxCol returns the index of the last column of a row given a cell index.
func (b *Board) GetRowMaxCol(i int) int {
	row := i / BoardHeight
	maxCol := row*BoardWidth + (BoardWidth - 1)
	return maxCol
}

// SetCells sets cells on the board.
func (b *Board) SetCells(cells []*Cell) Board {
	for _, c := range cells {
		if cell := b.GetAt(c.Index); cell != nil {
			b.SetAt(c.Index, *NewCell(&c.Tile, c.Index))
		} else {
			b.Cells = append(b.Cells, *c)
		}
	}
	sort.Slice(b.Cells, func(i, j int) bool {
		return b.Cells[i].Index < b.Cells[j].Index
	})
	return *b
}

// ReverseCells reverses cells.
func ReverseCells(cells []*Cell) []*Cell {
	reversed := []*Cell{}
	for c := len(cells) - 1; c >= 0; c-- {
		reversed = append(reversed, cells[c])
	}
	return reversed
}

// MarshalJSON serializes Tiles as a list of strings.
func (c Cell) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Tile  string `json:"t"`
		Index int    `json:"i"`
	}{
		Tile:  c.Tile.String(),
		Index: c.Index,
	})
}

// Print prints the board.
func (b *Board) Print() {
	dummy := NewBoard()
	dummyCells := []*Cell{}
	idx := 0
	for r := 0; r < BoardHeight; r++ {
		for c := 0; c < BoardWidth; c++ {
			cell := &Cell{Index: idx}
			dummyCells = append(dummyCells, cell)
			idx++
		}
	}
	dummy.SetCells(dummyCells)

	for _, c := range b.Cells {
		dummy.SetCells([]*Cell{&c})
	}

	printIdx := 0
	for r := 0; r < BoardHeight; r++ {
		for c := 0; c < BoardWidth; c++ {
			cell := dummy.GetAt(printIdx)
			if cell.Tile.Letter == "" {
				fmt.Printf("_ ")
			} else {
				fmt.Printf("%s ", cell.Tile.Letter)
			}
			printIdx++
		}
		fmt.Println()
	}
	fmt.Println()
}
