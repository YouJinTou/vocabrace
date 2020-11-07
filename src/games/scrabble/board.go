package scrabble

const _BoardWidth = 15
const _BoardHeight = 15

// Cell composes a board.
type Cell struct {
	tile                Tile
	row                 int
	col                 int
	isDoubleLetterScore bool
	isTripleLetterScore bool
	isDoubleWordScore   bool
	isTripleWordScore   bool
}

// Board is a 15x15 field of cells.
type Board struct {
	cells [_BoardHeight * _BoardWidth]Cell
}

// NewBoard creates a board.
func NewBoard() *Board {
	board := Board{
		cells: [_BoardHeight * _BoardWidth]Cell{},
	}

	for r := 0; r < _BoardHeight; r++ {
		for c := 0; c < _BoardWidth; c++ {
			cell := Cell{
				row: r,
				col: c,
			}
			board.cells[board.getCellIndex(r, c)] = cell
		}
	}

	return &board
}

// SetCell sets a tile at a particular cell.
func (b *Board) SetCell(r, c int, t Tile) Board {
	index := b.getCellIndex(r, c)
	b.cells[index].tile = t
	return *b
}

func (b *Board) getCellIndex(r, c int) int {
	index := (r * _BoardHeight) + c
	return index
}
