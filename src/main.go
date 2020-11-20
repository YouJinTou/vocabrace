package main

import (
	"encoding/json"
	"fmt"

	"github.com/YouJinTou/vocabrace/ws"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func main() {
	dc := scrabble.DynamoChecker{}
	dc.ValidateWords("bulgarian", []string{"полски"})
	bag := scrabble.NewBag("english")
	bt, _ := json.Marshal(&bag)
	nt := &scrabble.Bag{}
	err := json.Unmarshal(bt, nt)
	fmt.Println(err.Error())
	p := []*scrabble.Player{
		&scrabble.Player{
			ID:     "1",
			Name:   "Name",
			Points: 1,
		},
	}
	// ws.OnStart(&ws.ReceiverData{
	// 	Game:          "scrab5ble",
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
		PoolID:        "ac0c528e-0569-4140-a015-3a76c663bb8a",
		ConnectionIDs: []string{"a", "b"},
		Body:          "{\"Game\":\"scrabble\",\"IsPlace\":true,\"Word\":[{\"i\":1,\"t\":{\"l\":\"r\",\"v\":2}},{\"i\":2,\"t\":{\"l\":\"e\",\"v\":1}}]}",
	})
	x := scrabble.NewGame("english", p, scrabble.NewDynamoValidator())
	z := x.JSON()
	fmt.Println(z)
}
