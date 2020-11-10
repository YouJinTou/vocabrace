package ws

import (
	"encoding/json"
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

type scrabblews struct{}

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
			ToMove: game.ToMove.Name,
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
	type turn struct {
		IsPlace       bool
		IsExchange    bool
		IsPass        bool
		Word          string
		ExchangeTiles []string
	}
	current := turn{}
	err := json.Unmarshal([]byte(data.Body), &current)

	if err != nil {
		fmt.Println(fmt.Sprintf("Data: %s Dump: %s", data, err.Error()))
		return
	}

	game := &scrabble.Game{}
	loadState(data, game)

	if current.IsExchange {
		g, err := game.Exchange(current.ExchangeTiles)
		if err != nil {
			fmt.Println(fmt.Sprintf("Tile exchange. Data: %s Dump: %s", data, err.Error()))
			return
		}
		game = &g
	} else if current.IsPass {

	} else if current.IsPlace {

	} else {
		fmt.Println(fmt.Sprintf("Invalid action. Data: %s Dump: %s", data, err.Error()))
		return
	}

	if sErr := saveState(data, game); sErr != nil {
		panic(sErr.Error())
	}

	SendMany(data.otherConnections(), Message{
		Domain:  data.Domain,
		Stage:   data.Stage,
		Message: "Player placed word!",
	})
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
