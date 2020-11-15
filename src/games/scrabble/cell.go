package scrabble

// Cell composes a board.
type Cell struct {
	Tile  Tile `json:"t"`
	Index int  `json:"i"`
	m     *multiplier
}

// NewCell creates a new cell.
func NewCell(t *Tile, index int) *Cell {
	return &Cell{
		Tile:  *t,
		Index: index,
	}
}
