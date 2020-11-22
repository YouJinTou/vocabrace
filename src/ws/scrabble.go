package ws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

type scrabblews struct{}
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
type start struct {
	Tiles    *scrabble.Tiles `json:"t"`
	ToMove   string          `json:"m"`
	Players  []*player       `json:"p"`
	YourMove bool            `json:"y"`
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
type clientMessage struct {
	Body string
	Type string
}

func (s scrabblews) OnStart(data *ReceiverData) {
	players := s.loadPlayers(data.ConnectionIDs)
	game := scrabble.NewGame("bulgarian", players, scrabble.NewDynamoValidator())
	projected := s.projectPlayers(game.Players)
	messages := []*Message{}

	for _, p := range game.Players {
		startState := start{
			Tiles:    p.Tiles,
			ToMove:   game.ToMove().Name,
			Players:  projected,
			YourMove: game.ToMoveID == p.ID,
		}
		b, _ := json.Marshal(startState)
		messages = append(messages, &Message{
			Domain:       data.Domain,
			Stage:        data.Stage,
			ConnectionID: p.ID,
			Message:      string(b),
		})
	}

	if sErr := saveState(data, game); sErr != nil {
		panic(sErr.Error())
	}

	SendManyUnique(messages)
}

func (s scrabblews) OnAction(data *ReceiverData) {
	turn := turn{}
	bErr := json.Unmarshal([]byte(data.Body), &turn)

	if bErr != nil {
		s.returnClientError(data, "Invalid payload.", bErr)
		return
	}

	game := &scrabble.Game{}
	var r *result
	loadState(data, game)
	game.SetValidator(scrabble.NewDynamoValidator())

	if vErr := s.validateTurn(data, game); vErr != nil {
		s.returnClientError(data, "Invalid turn.", vErr)
		return
	}

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

	if sErr := saveState(data, game); sErr != nil {
		panic(sErr.Error())
	}

	Send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Domain,
		Stage:        data.Stage,
		Message:      r.d.JSONWithPersonal(),
	})

	messages := []*Message{}
	for _, cid := range data.otherConnections() {
		r.d.YourMove = game.ToMoveID == cid
		messages = append(messages, &Message{
			ConnectionID: cid,
			Domain:       data.Domain,
			Stage:        data.Stage,
			Message:      r.d.JSONWithoutPersonal(),
		})
	}
	SendManyUnique(messages)
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
		cells = append(cells, &scrabble.Cell{
			Tile:  scrabble.Tile{ID: c.TileID, Letter: c.Blank},
			Index: c.CellIndex,
		})
	}
	game, err := g.Place(scrabble.NewWord(cells))
	if err != nil {
		return &result{&game, scrabble.DeltaState{}, err}
	}
	return &result{&game, game.GetDelta(), err}
}

func (s *scrabblews) loadPlayers(connectionIDs []string) []*scrabble.Player {
	players := []*scrabble.Player{}
	for _, cid := range connectionIDs {
		players = append(players, &scrabble.Player{
			ID:     cid,
			Name:   cid,
			Points: 0,
		})
	}
	return players
}

func (s *scrabblews) projectPlayers(players []*scrabble.Player) []*player {
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

func (s *scrabblews) validateTurn(data *ReceiverData, g *scrabble.Game) error {
	if data.Initiator != g.ToMoveID {
		return errors.New("invalid player turn")
	}

	return nil
}

func (s *scrabblews) returnClientError(data *ReceiverData, message string, err error) {
	fmt.Println(fmt.Sprintf("ERROR Data: %s Dump: %s", data, err.Error()))
	msg := clientMessage{
		Body: message,
		Type: "ERROR",
	}
	b, _ := json.Marshal(msg)
	Send(&Message{
		ConnectionID: data.Initiator,
		Domain:       data.Domain,
		Stage:        data.Stage,
		Message:      string(b),
	})
}
