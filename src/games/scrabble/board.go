package scrabble

const _BoardWidth = 15
const _BoardHeight = 15

// Cell composes a board.
type Cell struct {
	Tile                Tile
	Row                 int
	Col                 int
	IsDoubleLetterScore bool
	IsTripleLetterScore bool
	IsDoubleWordScore   bool
	IsTripleWordScore   bool
}

// Board is a 15x15 field of cells.
type Board struct {
	Cells [_BoardHeight * _BoardWidth]Cell
}

// NewBoard creates a board.
func NewBoard() *Board {
	board := Board{
		Cells: [_BoardHeight * _BoardWidth]Cell{},
	}

	for r := 0; r < _BoardHeight; r++ {
		for c := 0; c < _BoardWidth; c++ {
			cell := Cell{
				Row: r,
				Col: c,
			}
			board.Cells[board.getCellIndex(r, c)] = cell
		}
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
