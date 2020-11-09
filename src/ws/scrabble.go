package ws

import (
	"encoding/json"
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func scrabbleOnStart(data *ReceiverData) {
	type start struct {
		Tiles  []*scrabble.Tile
		ToMove string
	}
	players := scrabbleLoadPlayers(data.ConnectionIDs)
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

	SendManyUnique(messages)
}

func scrabbleOnAction(data *ReceiverData) {
	type turn struct {
		IsPlace       bool
		IsExchange    bool
		IsPass        bool
		Word          string
		ExchangeCount int
	}
	current := turn{}
	err := json.Unmarshal([]byte(data.Body), &current)

	if err != nil {
		fmt.Println(fmt.Sprintf("Data: %s Dump: %s", data, err.Error()))
		return
	}

	if current.IsExchange {
		return
	}

	if current.IsPass {
		return
	}

	SendMany(data.otherConnections(), Message{
		Domain:  data.Domain,
		Stage:   data.Stage,
		Message: "Player placed word!",
	})
}

func scrabbleLoadPlayers(connectionIDs []string) []*scrabble.Player {
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
