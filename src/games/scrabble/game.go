package scrabble

import "encoding/json"

// Game holds a game's state.
type Game struct {
	Board *Board
	Bag   *Bag
}

// NewGame creates a new game.
func NewGame() *Game {
	return &Game{
		Board: NewBoard(),
		Bag:   NewBag(English),
	}
}

// JSON stringifies the game state to a JSON string.
func (g *Game) JSON() string {
	b, _ := json.Marshal(g)
	return string(b)
}
