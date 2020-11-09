package scrabble

// Player encapsulates player data.
type Player struct {
	ID     string  `json:"id"`
	Name   string  `json:"n"`
	Points int     `json:"p"`
	Tiles  []*Tile `json:"t"`
}
