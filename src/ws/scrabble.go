package ws

import (
	"encoding/json"

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
