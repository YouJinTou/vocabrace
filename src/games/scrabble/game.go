package scrabble

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
)

// Game holds a full game's state.
type Game struct {
	Board    Board     `json:"b"`
	Bag      Bag       `json:"g"`
	Players  []*Player `json:"p"`
	ToMoveID string    `json:"m"`
	Language string    `json:"l"`
	Order    []string  `json:"o"`
	delta    DeltaState
	v        CanValidate
}

// DeltaState shows the changes since the previous turn.
type DeltaState struct {
	ToMoveID             string      `json:"m"`
	LastAction           string      `json:"l"`
	LastActionPlayerID   string      `json:"i"`
	LastActionPlayerData interface{} `json:"d"`
	OtherPlayersData     interface{} `json:"o"`
	YourMove             bool        `json:"y"`
}

// JSONWithPersonal jsonifies a delta state with return data for the player.
func (d *DeltaState) JSONWithPersonal() string {
	p := DeltaState{
		ToMoveID:             d.ToMoveID,
		LastAction:           d.LastAction,
		LastActionPlayerID:   d.LastActionPlayerID,
		LastActionPlayerData: d.LastActionPlayerData,
		YourMove:             d.YourMove,
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
		YourMove:           d.YourMove,
	}
	b, _ := json.Marshal(p)
	result := string(b)
	return result
}

// NewGame creates a new game.
func NewGame(language string, players []*Player, validator CanValidate) *Game {
	if len(players) < 1 {
		panic("at least one player required")
	}

	bag := NewBag(language)
	for _, p := range players {
		p.Tiles = bag.Draw(7)
	}
	toMove, orderedIDs := orderPlayers(players)
	return &Game{
		Board:    *NewBoard(),
		Bag:      *bag,
		Players:  players,
		ToMoveID: toMove,
		Language: language,
		Order:    orderedIDs,
		v:        validator,
	}
}

// SetValidator sets the validator.
func (g *Game) SetValidator(v CanValidate) {
	g.v = v
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
func (g *Game) Place(w *Word) (Game, error) {
	if vErr := g.v.ValidatePlace(*g, w); vErr != nil {
		return *g, vErr
	}

	g.SetCellTiles(w.Cells)

	g.Board.SetCells(w.Cells)

	words := Extract(&g.Board, w)
	points := CalculatePoints(g, w, words)
	g.ToMove().AwardPoints(points)

	toReceive := g.Bag.Draw(w.Length())
	toRemove := []string{}
	for _, c := range w.Cells {
		toRemove = append(toRemove, c.Tile.ID)
	}
	_, err := g.ToMove().ExchangeTiles(toRemove, toReceive)

	g.delta = DeltaState{
		LastAction:           "Place",
		LastActionPlayerData: toReceive,
		OtherPlayersData:     w,
	}

	g.setNext()

	return *g, err
}

// SetCellTiles sets the incoming cells' tiles by looking them up by ID.
func (g *Game) SetCellTiles(cells []*Cell) error {
	for _, c := range cells {
		tile := g.ToMove().LookupTile(c.Tile.ID)

		if tile == nil {
			return fmt.Errorf("tile with ID %s not found", c.Tile.ID)
		}

		blank := c.Tile.Letter
		c.Tile = *tile
		if tile.IsBlank() {
			c.Tile.Letter = blank
		}
	}
	return nil
}

// FirstMovePlayed checks if the first move has been played.
func (g *Game) FirstMovePlayed() bool {
	return g.Board.GetAt(BoardOrigin) != nil
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
