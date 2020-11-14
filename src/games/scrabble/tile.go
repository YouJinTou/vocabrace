package scrabble

import "github.com/google/uuid"

// Tile represents a tile.
type Tile struct {
	Letter string `json:"l"`
	Value  int    `json:"v"`
	Index  string `json:"i"`
}

func tileID() string {
	return uuid.New().String()[0:5]
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value, tileID()}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
		Index:  tileID(),
	}
}

// Copy copies a tile.
func (t *Tile) Copy(preserveIndex bool) *Tile {
	var idx string
	if preserveIndex {
		idx = t.Index
	} else {
		idx = tileID()
	}
	return &Tile{
		Letter: t.Letter,
		Value:  t.Value,
		Index:  idx,
	}
}
