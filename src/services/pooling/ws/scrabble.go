package ws

import (
	"encoding/json"
	"errors"
	"log"

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
	PoolID   string          `json:"pid"`
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
	players := s.loadPlayers(data.Users)
	game := scrabble.NewGame(data.Language, players, scrabble.NewDynamoValidator())
	projected := s.setPlayerData(game.Players)
	messages := []*Message{}

	for _, p := range game.Players {
		startState := start{
			Tiles:    p.Tiles,
			ToMove:   game.ToMove().Name,
			Players:  projected,
			YourMove: game.ToMoveID == p.ID,
			PoolID:   data.PoolID,
		}
		b, _ := json.Marshal(startState)
		messages = append(messages, &Message{
			Domain:       data.Domain,
			Stage:        data.Stage,
			ConnectionID: userByID(data.Users, p.ID).ConnectionID,
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

func (s *scrabblews) validateTurn(data *ReceiverData, g *scrabble.Game) error {
	if data.Initiator != g.ToMoveID {
		return errors.New("invalid player turn")
	}

	return nil
}

func (s *scrabblews) returnClientError(data *ReceiverData, message string, err error) {
	log.Printf("ERROR Data: %+v Dump: %s", data, err.Error())
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
