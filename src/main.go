package main

import (
	"github.com/YouJinTou/vocabrace/ws"
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
		Initiator:     "WbTRaeF1liAAbIA=",
		Stage:         "dev",
		Domain:        "asd",
		PoolID:        "8e6eef18-33f0-44c6-92da-6c930299d45f",
		ConnectionIDs: []string{"a", "b"},
		Body:          "{\"g\":\"scrabble\",\"p\":true,\"w\":[{\"c\":110,\"t\":\"47d40a\",\"b\":null},{\"c\":111,\"t\":\"daa739\",\"b\":null},{\"c\":112,\"t\":\"1010f6\",\"b\":null},{\"c\":113,\"t\":\"53ed13\",\"b\":null}]}",
	})
	// x := scrabble.NewGame("english", p, scrabble.NewDynamoValidator())
	// z := x.JSON()
	// fmt.Println(z)
}
