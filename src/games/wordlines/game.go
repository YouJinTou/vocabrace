package wordlines

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"

	"github.com/YouJinTou/vocabrace/tools"
)

// Game holds a full game's state.
type Game struct {
	Board      Board     `json:"b"`
	Bag        Bag       `json:"g"`
	Players    []*Player `json:"p"`
	ToMoveID   string    `json:"m"`
	Language   string    `json:"l"`
	Order      []string  `json:"o"`
	EndCounter int       `json:"e"`
	delta      DeltaState
	v          CanValidate
}

// DeltaState shows the changes since the previous turn.
type DeltaState struct {
	ToMoveID             string         `json:"m"`
	LastAction           string         `json:"l"`
	LastActionPlayerID   string         `json:"i"`
	LastActionPlayerData interface{}    `json:"d"`
	OtherPlayersData     interface{}    `json:"o"`
	YourMove             bool           `json:"y"`
	Points               map[string]int `json:"p"`
	WinnerID             *string        `json:"w"`
	TilesRemaining       int            `json:"r"`
	Language             string         `json:"z"`
}

// JSONWithPersonal jsonifies a delta state with return data for the player.
func (d *DeltaState) JSONWithPersonal() string {
	p := DeltaState{
		ToMoveID:             d.ToMoveID,
		LastAction:           d.LastAction,
		LastActionPlayerID:   d.LastActionPlayerID,
		LastActionPlayerData: d.LastActionPlayerData,
		YourMove:             d.YourMove,
		Points:               d.Points,
		TilesRemaining:       d.TilesRemaining,
		Language:             d.Language,
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
		Points:             d.Points,
		TilesRemaining:     d.TilesRemaining,
		Language:           d.Language,
	}
	b, _ := json.Marshal(p)
	result := string(b)
	return result
}

// NewClassicGame creates a new classic game.
func NewClassicGame(language string, players []*Player, validator CanValidate) *Game {
	g := newGame(language, players, validator)
	g.Board = *NewBoard(classic{})
	return g
}

// NewSpiralGame creates a new spiral game.
func NewSpiralGame(language string, players []*Player, validator CanValidate) *Game {
	g := newGame(language, players, validator)
	g.Board = *NewBoard(spiral{})
	return g
}

func newGame(language string, players []*Player, validator CanValidate) *Game {
	if len(players) < 1 {
		panic("at least one player required")
	}

	bag := NewBag(language)
	for _, p := range players {
		p.Tiles = bag.Draw(7)
	}
	toMove, orderedIDs := orderPlayers(players)
	return &Game{
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
	g.EndCounter++

	g.handleEnd()

	return *g, err
}

// Pass passes a turn.
func (g *Game) Pass() Game {
	g.delta = DeltaState{LastAction: "Pass"}
	g.EndCounter++
	g.handleEnd()
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
	points := CalculatePoints(w, words, g.Board.l)
	g.ToMove().AwardPoints(points)

	toReceive := g.Bag.Draw(w.Length())
	toRemove := []string{}
	for _, c := range w.Cells {
		toRemove = append(toRemove, c.Tile.ID)
	}
	_, err := g.ToMove().ExchangeTiles(toRemove, toReceive)
	g.EndCounter = 0

	g.delta = DeltaState{
		LastAction:           "Place",
		LastActionPlayerData: toReceive,
		OtherPlayersData:     w,
	}

	g.handleEnd()

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

// IsVertical checks if a word is vertical on the board.
func (g *Game) IsVertical(w *Word) bool {
	if w.Length() == 1 {
		return true
	}
	indices := tools.SortInts(w.Indices()...)
	for i := 0; i < w.Length()-1; i++ {
		diff := indices[i+1] - indices[i]
		lettersSameCol := (diff % BoardHeight) == 0
		if !lettersSameCol {
			return false
		}
		areAdjacent := diff == BoardHeight
		if areAdjacent {
			continue
		}
		idx := indices[i] + BoardHeight
		for idx != indices[i+1] {
			if t := g.Board.GetAt(idx); t == nil {
				return false
			}
			idx += BoardHeight
		}
	}
	return true
}

// IsHorizontal checks if a word is horizontal on the board.
func (g *Game) IsHorizontal(w *Word) bool {
	if w.Length() == 1 {
		return true
	}
	indices := tools.SortInts(w.Indices()...)
	for i := 0; i < w.Length()-1; i++ {
		diff := indices[i+1] - indices[i]
		lettersSameRow := diff < BoardWidth
		if !lettersSameRow {
			return false
		}
		areAdjacent := diff == 1
		if areAdjacent {
			continue
		}
		idx := indices[i] + 1
		for idx != indices[i+1] {
			if t := g.Board.GetAt(idx); t == nil {
				return false
			}
			idx++
		}
	}
	return true
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

// GetPlayerByID gets a player by ID.
func (g *Game) GetPlayerByID(ID string) *Player {
	for _, p := range g.Players {
		if p.ID == ID {
			return p
		}
	}
	return nil
}

// GetLastMovedID returns the player ID that moved last.
func (g *Game) GetLastMovedID() string {
	return g.delta.LastActionPlayerID
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

// IsOver determines if a game is over.
func (g *Game) IsOver() bool {
	if g.IsPassiveEnd() {
		return true
	}
	if !g.Bag.IsEmpty() {
		return false
	}
	for _, p := range g.Players {
		if !p.HasTiles() {
			return true
		}
	}
	return false
}

// Leader gets the current points leader.
func (g *Game) Leader() *Player {
	max := -99999
	var leader *Player
	for _, p := range g.Players {
		if p.Points > max {
			leader = p
			max = p.Points
		}
	}
	return leader
}

// IsPassiveEnd checks if all players have exchanged or passed for the past n * 2 turns.
func (g *Game) IsPassiveEnd() bool {
	return g.EndCounter >= len(g.Players)*2
}

func (g *Game) playerPoints() map[string]int {
	m := map[string]int{}
	for _, p := range g.Players {
		m[p.ID] = p.Points
	}
	return m
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

func (g *Game) handleEnd() {
	if g.IsOver() {
		g.tallyFinalPoints()
		g.delta.WinnerID = &g.Leader().ID
	} else {
		g.setNext()
	}
	g.delta.Points = g.playerPoints()
	g.delta.TilesRemaining = g.Bag.Count()
}

func (g *Game) tallyFinalPoints() {
	if g.IsPassiveEnd() {
		for _, p := range g.Players {
			p.Points -= p.Tiles.Sum()
		}
		return
	}

	otherPlayersTilesSum := 0
	for _, p := range g.Players {
		if p.ID != g.ToMoveID {
			otherPlayersTilesSum += p.Tiles.Sum()
			p.Points -= p.Tiles.Sum()
		}
	}
	g.ToMove().Points += otherPlayersTilesSum
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
