package wordlines

import "github.com/YouJinTou/vocabrace/tools"

// Cell composes a board.
type Cell struct {
	Tile             Tile `json:"t"`
	Index            int  `json:"i"`
	enableMultiplier bool
}

// NewCell creates a new cell.
func NewCell(t *Tile, index int) *Cell {
	return &Cell{
		Tile:             *t,
		Index:            index,
		enableMultiplier: true,
	}
}

// Value calculates the cell's value given its tile value and any letter multipliers.
func (c *Cell) Value(l layout) int {
	if !c.enableMultiplier {
		return c.Tile.Value
	}
	return c.Tile.Value * c.LetterMultiplier(l)
}

// WordMultiplier returns the word multiplier of the cell.
func (c *Cell) WordMultiplier(l layout) int {
	if !c.enableMultiplier {
		return 1
	}
	if tools.ContainsInt(l.DoubleWordIndices(), c.Index) {
		return 2
	}
	if tools.ContainsInt(l.TripleWordIndices(), c.Index) {
		return 3
	}
	return 1
}

// LetterMultiplier returns the letter multiplier of the cell.
func (c *Cell) LetterMultiplier(l layout) int {
	if tools.ContainsInt(l.DoubleLetterIndices(), c.Index) {
		return 2
	}
	if tools.ContainsInt(l.TripleLetterIndices(), c.Index) {
		return 3
	}
	return 1
}

func cellTileIDs(cells []*Cell) []string {
	IDs := []string{}
	for _, c := range cells {
		IDs = append(IDs, c.Tile.ID)
	}
	return IDs
}
