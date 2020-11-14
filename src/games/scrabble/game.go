package scrabble

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

// Game holds a full game's state.
type Game struct {
	Board    *Board    `json:"b"`
	Bag      *Bag      `json:"g"`
	Players  []*Player `json:"p"`
	ToMoveID string    `json:"m"`
	Language string    `json:"l"`
	Order    []string  `json:"o"`
	delta    DeltaState
	v        *PlaceValidator
}

// DeltaState shows the changes since the previous turn.
type DeltaState struct {
	ToMoveID             string      `json:"m"`
	LastAction           string      `json:"l"`
	LastActionPlayerID   string      `json:"i"`
	LastActionPlayerData interface{} `json:"d"`
	OtherPlayersData     interface{} `json:"o"`
}

// JSONWithPersonal jsonifies a delta state with return data for the player.
func (d *DeltaState) JSONWithPersonal() string {
	p := DeltaState{
		ToMoveID:             d.ToMoveID,
		LastAction:           d.LastAction,
		LastActionPlayerID:   d.LastActionPlayerID,
		LastActionPlayerData: d.LastActionPlayerData,
	}
	b, _ := json.Marshal(p)
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
		v:        &PlaceValidator{},
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
func (g *Game) Exchange(ids []string) (Game, error) {
	if len(ids) == 0 {
		return *g, errors.New("must exchange at least one tile")
	}

	toReceive := g.Bag.Draw(len(ids))
	toReturn, err := g.ToMove().ExchangeTiles(ids, toReceive)
	g.Bag.Put(toReturn)

	g.delta = DeltaState{
		LastAction:           "Exchange",
		LastActionPlayerData: toReceive,
		OtherPlayersData:     toReceive.Count(),
	}

	g.setNext()

	return *g, err
}

// Pass passes a turn.
func (g *Game) Pass() Game {
	g.delta = DeltaState{
		LastAction: "Pass",
	}
	g.setNext()
	return *g
}

// Place places a word on the board.
func (g *Game) Place(tiles []*Cell) (Game, error) {
	if sErr := g.setCellTiles(tiles); sErr != nil {
		return *g, sErr
	}

	if vErr := g.v.ValidatePlaceAction(g, tiles); vErr != nil {
		return *g, vErr
	}

	g.Board.SetCells(tiles)

	points := calculatePlacePoints(g, tiles)
	g.ToMove().AwardPoints(points)

	toReceive := g.Bag.Draw(len(tiles))
	toRemove := []string{}
	for _, t := range tiles {
		toRemove = append(toRemove, t.Tile.ID)
	}
	_, err := g.ToMove().ExchangeTiles(toRemove, toReceive)

	g.delta = DeltaState{
		LastAction:           "Place",
		LastActionPlayerData: toReceive,
		OtherPlayersData:     tiles,
	}

	g.setNext()

	return *g, err
}

func (g *Game) setCellTiles(cells []*Cell) error {
	for _, c := range cells {
		tile := g.ToMove().LookupTile(c.Tile.ID)
		if tile == nil {
			return fmt.Errorf("tile with ID %s not found", c.Tile.ID)
		}
		c.Tile = *tile
	}
	return nil
}

func (g *Game) setNext() {
	for i, p := range g.Order {
		if p == g.ToMoveID {
			g.delta.LastActionPlayerID = g.ToMoveID
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

// LastToMove gets the player who moved last.
func (g *Game) LastToMove() *Player {
	for _, p := range g.Players {
		if p.ID == g.delta.LastActionPlayerID {
			return p
		}
	}

	panic("cannot find last player")
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
