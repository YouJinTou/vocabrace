package main

import (
	"github.com/YouJinTou/vocabrace/services/pooling/ws"
)

func main() {
	// dc := scrabble.DynamoChecker{}
	// dc.ValidateWords("bulgarian", []string{"полски", "таралежка", "пинокио", "глупак"})
	// bag := scrabble.NewBag("english")
	// bt, _ := json.Marshal(&bag)
	// nt := &scrabble.Bag{}
	// err := json.Unmarshal(bt, nt)
	// fmt.Println(err.Error())
	// p := []*scrabble.Player{
	// 	&scrabble.Player{
	// 		ID:     "1",
	// 		Name:   "Name",
	// 		Points: 1,
	// 	},
	// }
	// ws.OnStart(&ws.ReceiverData{
	// 	Game:          "scrab5ble",
	// 	Stage:         "dev",
	// 	Domain:        "asd",
	// 	PoolID:        "ac0c528e-0569-4140-a015-3a76c663bb8a",
	// 	ConnectionIDs: []string{"a", "b"},
	// })
	ws.OnAction(&ws.ReceiverData{
		Game:          "scrabble",
		Initiator:     "WbW2EdWqliACGZQ=",
		Stage:         "dev",
		Domain:        "asd",
		PoolID:        "c6c1b692-5743-4950-8229-f03387c6386e",
		ConnectionIDs: []string{"a", "b"},
		Body:          "{\"g\":\"scrabble\",\"p\":true,\"w\":[{\"c\":110,\"t\":\"1b283a\",\"b\":null},{\"c\":111,\"t\":\"d49a35\",\"b\":null},{\"c\":112,\"t\":\"f7c6a6\",\"b\":null},{\"c\":113,\"t\":\"70de9b\",\"b\":null}]}",
	})
	// x := scrabble.NewGame("english", p, scrabble.NewDynamoValidator())
	// z := x.JSON()
	// fmt.Println(z)
}
