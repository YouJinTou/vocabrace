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
	ws.OnStart(&ws.ReceiverData{
		Game:          "scrabble",
		Stage:         "dev",
		Domain:        "asd",
		PoolID:        "ac0c528e-0569-4140-a015-3a76c663bb8a",
		ConnectionIDs: []string{"a", "b"},
	})
	x := scrabble.NewGame(p)
	z := x.JSON()
	fmt.Println(z)
}
