package scrabble

import (
	"encoding/json"
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

// Cell composes a board.
type Cell struct {
	Tile  Tile `json:"t"`
	Index int  `json:"i"`
}

// NewCell creates a new cell.
func NewCell(t *Tile, index int) *Cell {
	return &Cell{
		Tile:  *t,
		Index: index,
	}
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
