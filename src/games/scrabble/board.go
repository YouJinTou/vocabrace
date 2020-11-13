package scrabble

import "sort"

const _BoardWidth = 15
const _BoardHeight = 15

// Cell composes a board.
type Cell struct {
	Tile  Tile `json:"t"`
	Index int  `json:"i"`
}

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

// SetCells sets tiles on the board.
func (b *Board) SetCells(cells []*Cell) Board {
	for _, c := range cells {
		b.Cells = append(b.Cells, *c)
	}
	sort.Slice(b.Cells, func(i, j int) bool {
		return b.Cells[i].Index < b.Cells[j].Index
	})
	return *b
}
