package wordlines

import (
	"strings"
)

// Word encapsulates a valid arrangement of tiles.
type Word struct {
	Cells []*Cell
}

// NewWord creates a new word.
func NewWord(cells []*Cell) *Word {
	if len(cells) == 0 {
		panic("word must not be empty")
	}
	return &Word{Cells: cells}
}

// NewWordFromString creates a word given a string and its indices on the board.
func NewWordFromString(s string, values, indices []int) *Word {
	tokens := strings.Split(s, "")
	cells := []*Cell{}
	for i, idx := range indices {
		tile := NewTile(tokens[i], values[i])
		cells = append(cells, NewCell(tile, idx))
	}
	w := NewWord(cells)
	return w
}

// ToStrings converts words to their string representations.
func ToStrings(words []*Word) []string {
	result := []string{}
	for _, w := range words {
		result = append(result, w.String())
	}
	return result
}

// String returns the string representation of a word.
func (w *Word) String() string {
	var s string
	for _, c := range w.Cells {
		s += c.Tile.Letter
	}
	return s
}

// Length gets a word's length.
func (w *Word) Length() int {
	return len(w.Cells)
}

// IsSameAs checks if two words match.
func (w *Word) IsSameAs(o *Word) bool {
	if w.String() != o.String() {
		return false
	}

	for i, c := range w.Cells {
		if c.Index != o.Cells[i].Index {
			return false
		}
	}

	return true
}

// ExistsIn checks if a word exists in a set of words.
func (w *Word) ExistsIn(words []*Word) bool {
	for _, word := range words {
		if w.IsSameAs(word) {
			return true
		}
	}
	return false
}

// Value returns the sum of its tiles.
func (w *Word) Value(l layout) int {
	sum := 0
	for _, c := range w.Cells {
		sum += c.Value(l)
	}
	for _, c := range w.Cells {
		sum *= c.WordMultiplier(l)
	}
	return sum
}

// Extract returns the current word plus any adjacently formed words.
func Extract(b *Board, w *Word) []*Word {
	words := []*Word{}
	for _, c := range w.Cells {
		if w := traverseVertically(b, c); w != nil {
			words = append(words, w)
		}
		if w := traverseHorizontally(b, c); w != nil {
			words = append(words, w)
		}
	}

	result := removeUnnecessary(w, words)
	return result
}

// Indices returns a list of a word's cell indices.
func (w *Word) Indices() []int {
	result := []int{}
	for _, c := range w.Cells {
		result = append(result, c.Index)
	}
	return result
}

func traverseVertically(b *Board, c *Cell) *Word {
	cells := []*Cell{}
	for i := c.Index - BoardHeight; i >= BoardMinIndex; i -= BoardHeight {
		if c := b.GetAt(i); c != nil {
			cells = append(cells, c)
		} else {
			break
		}
	}
	cells = append(ReverseCells(cells), c)
	for i := c.Index + BoardHeight; i <= BoardMaxIndex; i += BoardHeight {
		if c := b.GetAt(i); c != nil {
			cells = append(cells, c)
		} else {
			break
		}
	}

	return NewWord(cells)
}

func traverseHorizontally(b *Board, c *Cell) *Word {
	cells := []*Cell{}
	for i := c.Index - 1; i >= b.GetRowMinCol(c.Index); i-- {
		if c := b.GetAt(i); c != nil {
			cells = append(cells, c)
		} else {
			break
		}
	}
	cells = append(ReverseCells(cells), c)
	for i := c.Index + 1; i <= b.GetRowMaxCol(c.Index); i++ {
		if c := b.GetAt(i); c != nil {
			cells = append(cells, c)
		} else {
			break
		}
	}

	return NewWord(cells)
}

func removeUnnecessary(initial *Word, words []*Word) []*Word {
	result := []*Word{}
	for _, w := range words {
		isSingleLetter := w.Length() == 1

		if !isSingleLetter && !w.ExistsIn(result) {
			result = append(result, w)
		}
	}
	return result
}
