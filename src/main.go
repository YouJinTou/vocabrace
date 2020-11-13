package main

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/ws"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func main() {
	p := []*scrabble.Player{
		&scrabble.Player{
			ID:     "1",
			Name:   "Name",
			Points: 1,
		},
	}
	// ws.OnStart(&ws.ReceiverData{
	// 	Game:          "scrabble",
	// 	Stage:         "dev",
	// 	Domain:        "asd",
	// 	PoolID:        "ac0c528e-0569-4140-a015-3a76c663bb8a",
	// 	ConnectionIDs: []string{"a", "b"},
	// })
	ws.OnAction(&ws.ReceiverData{
		Game:          "scrabble",
		Initiator:     "V9MXTex9FiACFew=",
		Stage:         "dev",
		Domain:        "asd",
		PoolID:        "6e236582-1b8b-43ab-bce2-64dbbf733e0b",
		ConnectionIDs: []string{"a", "b"},
		Body:          "{\"Game\":\"scrabble\",\"IsPlace\":true,\"Word\":[{\"i\":1,\"t\":{\"l\":\"r\",\"v\":2}},{\"i\":2,\"t\":{\"l\":\"e\",\"v\":1}}]}",
	})
	x := scrabble.NewGame(p)
	z := x.JSON()
	fmt.Println(z)
}
