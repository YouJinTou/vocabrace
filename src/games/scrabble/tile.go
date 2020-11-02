package scrabble

// Tile represents a tile.
type Tile struct {
	Letter string
	Value  int
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
