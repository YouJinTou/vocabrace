package scrabble

import (
	"fmt"
)

// Game holds a game's state.
type Game struct {
	board *Board
	bag   *Bag
}

// NewGame creates a new game.
func NewGame() *Game {
	return &Game{
		board: NewBoard(),
		bag:   NewBag(English),
	}
}

func (g *Game) Print() {
	fmt.Println("HEEE")
}
