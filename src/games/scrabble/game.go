package scrabble

import (
	"encoding/json"
	"math/rand"
	"strconv"
)

// Game holds a full game's state.
type Game struct {
	Board    *Board    `json:"b"`
	Bag      *Bag      `json:"bg"`
	Players  []*Player `json:"p"`
	ToMoveID string    `json:"m"`
	Language string    `json:"l"`
	Order    []string  `json:"o"`
	delta    DeltaState
}

// DeltaState shows the changes since the previous turn.
type DeltaState struct {
	ToMoveID             string
	LastAction           string
	LastActionPlayerID   string
	LastActionPlayerData string
	OtherPlayersData     string
}

// JSONWithPersonal jsonifies a delta state with return data for the player.
func (d *DeltaState) JSONWithPersonal() string {
	b, _ := json.Marshal(d)
	result := string(b)
	return result
}

// JSONWithoutPersonal jsonifies a delta state without return data for the player.
func (d *DeltaState) JSONWithoutPersonal() string {
	p := DeltaState{
		ToMoveID:           d.ToMoveID,
		LastAction:         d.LastAction,
		LastActionPlayerID: d.LastActionPlayerID,
		OtherPlayersData:   d.OtherPlayersData,
	}
	b, _ := json.Marshal(p)
	result := string(b)
	return result
}

// NewGame creates a new game.
func NewGame(players []*Player) *Game {
	if len(players) < 1 {
		panic("at least one player required")
	}

	bag := NewBag(English)
	for _, p := range players {
		p.Tiles = bag.Draw(7)
	}
	toMove, orderedIDs := orderPlayers(players)
	return &Game{
		Board:    NewBoard(),
		Bag:      bag,
		Players:  players,
		ToMoveID: toMove,
		Language: "en",
		Order:    orderedIDs,
	}
}

// GetDelta shows the changes since the last move.
func (g *Game) GetDelta() DeltaState {
	return g.delta
}

// JSON stringifies the game state to a JSON string.
func (g *Game) JSON() string {
	b, _ := json.Marshal(g)
	return string(b)
}

// Exchange exchanges a set of tiles for the player to move.
func (g *Game) Exchange(exchangeTiles []string) (Game, error) {
	toReceive := g.Bag.Draw(len(exchangeTiles))
	toReceiveBytes, _ := json.Marshal(toReceive)
	toReturn, err := g.ToMove().ExchangeTiles(exchangeTiles, toReceive)
	g.Bag.Put(toReturn)

	g.delta = DeltaState{
		LastAction:           "EXCHANGE",
		LastActionPlayerID:   g.ToMove().ID,
		LastActionPlayerData: string(toReceiveBytes),
		OtherPlayersData:     strconv.Itoa(len(toReceive)),
	}

	g.setNext()

	return *g, err
}

func (g *Game) setNext() {
	for i, p := range g.Order {
		if p == g.ToMoveID {
			if i+1 == len(g.Order) {
				g.ToMoveID = g.Order[0]
			} else {
				g.ToMoveID = g.Order[i+1]
			}
			g.delta.ToMoveID = g.ToMoveID
			return
		}
	}
}

// ToMove gets the player to move.
func (g *Game) ToMove() *Player {
	for _, p := range g.Players {
		if p.ID == g.ToMoveID {
			return p
		}
	}

	panic("cannot find player")
}

func orderPlayers(players []*Player) (string, []string) {
	toMoveIdx := rand.Intn(len(players))
	toMove := players[toMoveIdx]
	orderedPlayers := []*Player{toMove}
	orderedPlayers = append(orderedPlayers, players[toMoveIdx+1:]...)
	orderedPlayers = append(orderedPlayers, players[0:toMoveIdx]...)
	ids := []string{}
	for _, p := range orderedPlayers {
		ids = append(ids, p.ID)
	}
	return orderedPlayers[0].ID, ids
}
