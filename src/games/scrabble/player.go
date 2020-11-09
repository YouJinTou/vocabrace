package scrabble

// Player encapsulates player data.
type Player struct {
	ID     string
	Name   string
	Points int
	Tiles  []*Tile
}
