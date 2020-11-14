package scrabble

import "github.com/google/uuid"

// Tile represents a tile.
type Tile struct {
	Letter string `json:"l"`
	Value  int    `json:"v"`
	index  string
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value, uuid.New().String()}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
		index:  uuid.New().String(),
	}
}

// Copy copies a tile.
func (t *Tile) Copy(preserveIndex bool) *Tile {
	var idx string
	if preserveIndex {
		idx = t.index
	} else {
		idx = uuid.New().String()
	}
	return &Tile{
		Letter: t.Letter,
		Value:  t.Value,
		index:  idx,
	}
}
