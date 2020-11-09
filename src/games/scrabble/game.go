package scrabble

import (
	"encoding/json"
	"math/rand"
)

// Game holds a full game's state.
type Game struct {
	Board   *Board
	Bag     *Bag
	Players []*Player
	ToMove  *Player
}

// NewGame creates a new game.
func NewGame(players []*Player) *Game {
	bag := NewBag(English)
	for _, p := range players {
		p.Tiles = bag.Draw(7)
	}
	return &Game{
		Board:   NewBoard(),
		Bag:     bag,
		Players: players,
		ToMove:  players[rand.Intn(len(players))],
	}
}

// JSON stringifies the game state to a JSON string.
func (g *Game) JSON() string {
	b, _ := json.Marshal(g)
	return string(b)
}
