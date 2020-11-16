package scrabble

import "github.com/YouJinTou/vocabrace/tools"

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

// Value calculates the cell's value given its tile value and any letter multipliers.
func (c *Cell) Value() int {
	return c.Tile.Value * multiplier(c.Index)
}

// WordMultiplier returns the word multiplier of the cell.
func (c *Cell) WordMultiplier() int {
	doubleWordIndices := []int{16, 32, 48, 64, BoardOrigin, 160, 176, 192, 208}
	if tools.ContainsInt(doubleWordIndices, c.Index) {
		return 2
	}
	tripleWordIndices := []int{0, 7, 14, 105, 119, 210, 224}
	if tools.ContainsInt(tripleWordIndices, c.Index) {
		return 3
	}
	return 1
}

func multiplier(i int) int {
	doubleLetterIndices := []int{3, 11, 36, 38, 45, 52, 59, 92, 96, 98, 102, 108, 116, 122, 126, 128, 132, 165, 172, 179, 186, 188, 213, 221}
	if tools.ContainsInt(doubleLetterIndices, i) {
		return 2
	}
	tripleLetterIndices := []int{20, 24, 76, 80, 84, 88, 136, 140, 144, 148, 200, 204}
	if tools.ContainsInt(tripleLetterIndices, i) {
		return 3
	}
	return 1
}
