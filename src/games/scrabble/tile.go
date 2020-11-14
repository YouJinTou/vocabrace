package scrabble

import "github.com/google/uuid"

// Tile represents a tile.
type Tile struct {
	Letter string `json:"l"`
	Value  int    `json:"v"`
	Index  string `json:"i"`
}

func tileIndex() string {
	return uuid.New().String()[0:5]
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value, tileIndex()}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
		Index:  tileIndex(),
	}
}

// Copy copies a tile.
func (t *Tile) Copy(preserveIndex bool) *Tile {
	var idx string
	if preserveIndex {
		idx = t.Index
	} else {
		idx = tileIndex()
	}
	return &Tile{
		Letter: t.Letter,
		Value:  t.Value,
		Index:  idx,
	}
}
