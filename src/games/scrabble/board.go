package scrabble

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

// SetCell sets a tile at a particular cell.
func (b *Board) SetCell(r, c int, t Tile) Board {
	index := b.getCellIndex(r, c)
	b.Cells[index].Tile = t
	return *b
}

func (b *Board) getCellIndex(r, c int) int {
	index := (r * _BoardHeight) + c
	return index
}
