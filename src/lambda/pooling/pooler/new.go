package main

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func loadInitialState(game string, connectionIDs []string) string {
	players := loadPlayers(connectionIDs)
	switch game {
	case "scrabble":
		return loadScrabble(players)
	default:
		panic(fmt.Sprintf("could not load initial state for %s", game))
	}
}

func loadPlayers(connectionIDs []string) []*scrabble.Player {
	players := []*scrabble.Player{}
	for _, cid := range connectionIDs {
		players = append(players, &scrabble.Player{
			Name:   cid,
			Points: 0,
		})
	}
	return players
}

func loadScrabble(players []*scrabble.Player) string {
	game := scrabble.NewGame(players)
	return game.JSON()
}
