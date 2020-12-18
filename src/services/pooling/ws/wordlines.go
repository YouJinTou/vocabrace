package ws

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/YouJinTou/vocabrace/games/wordlines"
	"github.com/google/uuid"
)

type wordlinesws struct {
	loadState      func(string, interface{})
	saveState      func(*saveStateInput) error
	send           func(*Message) error
	sendManyUnique func([]*Message)
}
type cell struct {
	CellIndex int    `json:"c"`
	TileID    string `json:"t"`
	Blank     string `json:"b"`
}
type player struct {
	ID     string `json:"i"`
	Name   string `json:"n"`
	Points int    `json:"p"`
	Tiles  int    `json:"t"`
}
type turn struct {
	IsPlace       bool     `json:"p"`
	IsExchange    bool     `json:"e"`
	IsPass        bool     `json:"q"`
	Word          []*cell  `json:"w"`
	ExchangeTiles []string `json:"t"`
}
type result struct {
	g   *wordlines.Game
	d   wordlines.DeltaState
	err error
}

func (s wordlinesws) OnStart(c *Connections) {
	players := s.loadPlayers(c)
	game := wordlines.NewSpiralGame(c.Language(), players, wordlines.NewDynamoValidator())
	projected := s.setPlayerData(game.Players)
	messages := []*Message{}
	startState := struct {
		Tiles          *wordlines.Tiles `json:"t"`
		ToMoveID       string           `json:"m"`
		Players        []*player        `json:"p"`
		YourMove       bool             `json:"y"`
		PoolID         string           `json:"pid"`
		TilesRemaining int              `json:"r"`
		Language       string           `json:"z"`
	}{nil, game.ToMoveID, projected, false, uuid.New().String(), game.Bag.Count(), c.Language()}

	for _, p := range game.Players {
		startState.Tiles = p.Tiles
		startState.YourMove = game.ToMoveID == p.ID
		b, _ := json.Marshal(startState)
		messages = append(messages, &Message{
			Domain:       c.Domain(),
			ConnectionID: *c.IDByUserID(p.ID),
			Message:      string(b),
		})
	}

	if sErr := s.saveState(&saveStateInput{
		PoolID:        startState.PoolID,
		ConnectionIDs: c.IDs(),
		V:             game,
	}); sErr != nil {
		panic(sErr.Error())
	}

	s.sendManyUnique(messages)
}

func (s wordlinesws) OnAction(data *OnActionInput) {
	turn := turn{}
	bErr := json.Unmarshal([]byte(data.Body), &turn)

	if bErr != nil {
		s.returnClientError(data, "Invalid payload.", bErr)
		return
	}

	game := &wordlines.Game{}
	s.loadState(data.PoolID, game)
	game.SetValidator(wordlines.NewDynamoValidator())

	if vErr := s.validateTurn(data, game); vErr != nil {
		s.returnClientError(data, "Invalid turn.", vErr)
		return
	}

	var r *result
	if turn.IsExchange {
		r = s.exchange(&turn, game)
	} else if turn.IsPass {
		r = s.pass(game)
	} else if turn.IsPlace {
		r = s.place(&turn, game)
	} else {
		r = &result{nil, wordlines.DeltaState{}, errors.New("invalid action")}
	}

	if r.err != nil {
		s.returnClientError(data, "Bad request.", r.err)
		return
	}

	if sErr := s.saveState(&saveStateInput{
		PoolID:        data.PoolID,
		ConnectionIDs: data.Connections.IDs(),
		V:             game,
	}); sErr != nil {
		panic(sErr.Error())
	}

	s.send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Connections.Domain(),
		Message:      r.d.JSONWithPersonal(),
	})

	messages := []*Message{}
	for _, cid := range data.Connections.OtherIDs(data.Initiator) {
		r.d.YourMove = game.ToMoveID == *data.Connections.UserIDByID(cid)
		messages = append(messages, &Message{
			ConnectionID: cid,
			Domain:       data.Connections.Domain(),
			Message:      r.d.JSONWithoutPersonal(),
		})
	}
	s.sendManyUnique(messages)
}

func (s *wordlinesws) exchange(turn *turn, g *wordlines.Game) *result {
	game, err := g.Exchange(turn.ExchangeTiles)
	if err != nil {
		return &result{&game, wordlines.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s *wordlinesws) pass(g *wordlines.Game) *result {
	game := g.Pass()
	return &result{&game, game.GetDelta(), nil}
}

func (s *wordlinesws) place(turn *turn, g *wordlines.Game) *result {
	cells := []*wordlines.Cell{}
	for _, c := range turn.Word {
		cells = append(cells, wordlines.NewCell(
			wordlines.NewTileWithID(c.TileID, c.Blank, 0),
			c.CellIndex,
		))
	}
	game, err := g.Place(wordlines.NewWord(cells))
	if err != nil {
		return &result{&game, wordlines.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s *wordlinesws) loadPlayers(c *Connections) []*wordlines.Player {
	result := []*wordlines.Player{}
	for _, u := range c.UserIDs() {
		result = append(result, &wordlines.Player{
			ID:     u,
			Name:   u,
			Points: 0,
		})
	}
	return result
}

func (s *wordlinesws) setPlayerData(players []*wordlines.Player) []*player {
	result := []*player{}
	for _, p := range players {
		result = append(result, &player{
			ID:     p.ID,
			Name:   p.Name,
			Points: p.Points,
			Tiles:  p.Tiles.Count(),
		})
	}
	return result
}

func (s *wordlinesws) validateTurn(data *OnActionInput, g *wordlines.Game) error {
	if data.InitiatorUserID != g.ToMoveID {
		return errors.New("invalid player turn")
	}

	return nil
}

func (s *wordlinesws) returnClientError(data *OnActionInput, message string, err error) {
	log.Printf("ERROR Data: %+v Dump: %s", data, err.Error())
	msg := struct {
		Body string
		Type string
	}{message, "ERROR"}
	b, _ := json.Marshal(msg)
	s.send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Connections.Domain(),
		Message:      string(b),
	})
}
