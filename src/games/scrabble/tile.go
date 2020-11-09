package scrabble

// Tile represents a tile.
type Tile struct {
	Letter string `json:"l"`
	Value  int    `json:"v"`
}

// NewTile creates a new tile.
func NewTile(letter string, value int) *Tile {
	return &Tile{letter, value}
}

// BlankTile creates a blank tile.
func BlankTile() *Tile {
	return &Tile{
		Letter: "",
		Value:  0,
	}
}

// Copy copies a tile.
func (t *Tile) Copy() *Tile {
	return &Tile{
		Letter: t.Letter,
		Value:  t.Value,
	}
}
