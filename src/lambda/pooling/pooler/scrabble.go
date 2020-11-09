package main

import (
	"encoding/json"

	"github.com/YouJinTou/vocabrace/games/scrabble"
	ws "github.com/YouJinTou/vocabrace/lambda/pooling"
)

func onStart(pool *pool) {
	type start struct {
		Tiles  []*scrabble.Tile
		ToMove string
	}
	players := loadPlayers(pool.ConnectionIDs())
	game := scrabble.NewGame(players)
	messages := []*ws.Message{}
	for _, p := range game.Players {
		startState := start{
			Tiles:  p.Tiles,
			ToMove: game.ToMove.Name,
		}
		b, _ := json.Marshal(startState)
		messages = append(messages, &ws.Message{
			Domain:       pool.Domain,
			Stage:        pool.Stage,
			ConnectionID: p.ID,
			Message:      string(b),
		})
	}
	ws.SendManyUnique(messages)
}

func loadPlayers(connectionIDs []string) []*scrabble.Player {
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
