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

// DeltaState shows the changes since the previous turn.
type DeltaState struct {
	ToMoveID   string
	LastAction string
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

// GetDelta shows the changes since the last move.
func (g *Game) GetDelta() DeltaState {
	return DeltaState{}
}

// JSON stringifies the game state to a JSON string.
func (g *Game) JSON() string {
	b, _ := json.Marshal(g)
	return string(b)
}
