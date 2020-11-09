package scrabble

// Player encapsulates player data.
type Player struct {
	Name   string
	Points int
	Tiles  []*Tile
}
