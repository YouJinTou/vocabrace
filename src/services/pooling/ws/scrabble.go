package ws

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

type scrabblews struct {
	loadState      func(string, interface{})
	saveState      func(string, interface{}) error
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
	g   *scrabble.Game
	d   scrabble.DeltaState
	err error
}

func (s scrabblews) OnStart(data *OnStartInput) {
	players := s.loadPlayers(data.Users)
	game := scrabble.NewGame(data.Language, players, scrabble.NewDynamoValidator())
	projected := s.setPlayerData(game.Players)
	messages := []*Message{}
	startState := struct {
		Tiles    *scrabble.Tiles `json:"t"`
		ToMove   string          `json:"m"`
		Players  []*player       `json:"p"`
		YourMove bool            `json:"y"`
		PoolID   string          `json:"pid"`
	}{nil, game.ToMove().Name, projected, false, data.PoolID}

	for _, p := range game.Players {
		startState.Tiles = p.Tiles
		startState.YourMove = game.ToMoveID == p.ID
		b, _ := json.Marshal(startState)
		messages = append(messages, &Message{
			Domain:       data.Domain,
			ConnectionID: userByID(data.Users, p.ID).ConnectionID,
			Message:      string(b),
		})
	}

	if sErr := s.saveState(data.PoolID, game); sErr != nil {
		panic(sErr.Error())
	}

	s.sendManyUnique(messages)
}

func (s scrabblews) OnAction(data *OnActionInput) {
	turn := turn{}
	bErr := json.Unmarshal([]byte(data.Body), &turn)

	if bErr != nil {
		s.returnClientError(data, "Invalid payload.", bErr)
		return
	}

	game := &scrabble.Game{}
	s.loadState(data.PoolID, game)
	game.SetValidator(scrabble.NewDynamoValidator())

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
		r = &result{nil, scrabble.DeltaState{}, errors.New("invalid action")}
	}

	if r.err != nil {
		s.returnClientError(data, "Bad request.", r.err)
		return
	}

	if sErr := s.saveState(data.PoolID, game); sErr != nil {
		panic(sErr.Error())
	}

	s.send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Domain,
		Message:      r.d.JSONWithPersonal(),
	})

	messages := []*Message{}
	for _, cid := range data.otherConnections() {
		r.d.YourMove = game.ToMoveID == cid
		messages = append(messages, &Message{
			ConnectionID: cid,
			Domain:       data.Domain,
			Message:      r.d.JSONWithoutPersonal(),
		})
	}
	s.sendManyUnique(messages)
}

func (s *scrabblews) exchange(turn *turn, g *scrabble.Game) *result {
	game, err := g.Exchange(turn.ExchangeTiles)
	if err != nil {
		return &result{&game, scrabble.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s *scrabblews) pass(g *scrabble.Game) *result {
	game := g.Pass()
	return &result{&game, game.GetDelta(), nil}
}

func (s *scrabblews) place(turn *turn, g *scrabble.Game) *result {
	cells := []*scrabble.Cell{}
	for _, c := range turn.Word {
		cells = append(cells, scrabble.NewCell(
			scrabble.NewTileWithID(c.TileID, c.Blank, 0),
			c.CellIndex,
		))
	}
	game, err := g.Place(scrabble.NewWord(cells))
	if err != nil {
		return &result{&game, scrabble.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s *scrabblews) loadPlayers(users []*User) []*scrabble.Player {
	result := []*scrabble.Player{}
	for _, u := range users {
		result = append(result, &scrabble.Player{
			ID:     u.UserID,
			Name:   u.Username,
			Points: 0,
		})
	}
	return result
}

func (s *scrabblews) setPlayerData(players []*scrabble.Player) []*player {
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

func (s *scrabblews) validateTurn(data *OnActionInput, g *scrabble.Game) error {
	if data.Initiator != g.ToMoveID {
		return errors.New("invalid player turn")
	}

	return nil
}

func (s *scrabblews) returnClientError(data *OnActionInput, message string, err error) {
	log.Printf("ERROR Data: %+v Dump: %s", data, err.Error())
	msg := struct {
		Body string
		Type string
	}{message, "ERROR"}
	b, _ := json.Marshal(msg)
	s.send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Domain,
		Message:      string(b),
	})
}
