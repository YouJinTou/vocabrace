package ws

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

type scrabblews struct{}
type turn struct {
	IsPlace       bool
	IsExchange    bool
	IsPass        bool
	Word          string
	ExchangeTiles []string
}
type result struct {
	g   *scrabble.Game
	m   string
	err error
}
type clientMessage struct {
	Body string
	Type string
}

func (s scrabblews) OnStart(data *ReceiverData) {
	type start struct {
		Tiles  []*scrabble.Tile
		ToMove string
	}
	players := s.loadPlayers(data.ConnectionIDs)
	game := scrabble.NewGame(players)
	messages := []*Message{}

	for _, p := range game.Players {
		startState := start{
			Tiles:  p.Tiles,
			ToMove: game.ToMove().Name,
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

	if vErr := s.validateTurn(data, game); vErr != nil {
		s.returnClientError(data, "Invalid turn.", vErr)
		return
	}

	if turn.IsExchange {
		r = s.exchange(&turn, game)
	} else if turn.IsPass {

	} else if turn.IsPlace {

	} else {
		r = &result{nil, "", errors.New("invalid action")}
	}

	if r.err != nil {
		s.returnClientError(data, "Bad request.", r.err)
		return
	}

	if sErr := saveState(data, game); sErr != nil {
		panic(sErr.Error())
	}

	SendMany(data.otherConnections(), Message{
		Domain:  data.Domain,
		Stage:   data.Stage,
		Message: r.m,
	})
}

func (s *scrabblews) exchange(turn *turn, g *scrabble.Game) *result {
	game, err := g.Exchange(turn.ExchangeTiles)
	if err != nil {
		return &result{&game, "EXCHANGE FAILED", err}
	}
	return &result{&game, game.GetDeltaJSON(), err}
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
