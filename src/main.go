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
		Initiator:     "WbQJAeMaFiACGJA=",
		Stage:         "dev",
		Domain:        "asd",
		PoolID:        "08b54d5a-91e3-45be-a2c1-d65ba389fb14",
		ConnectionIDs: []string{"a", "b"},
		Body:          "{\"g\":\"scrabble\",\"p\":true,\"w\":[{\"c\":112,\"t\":\"e17a43\",\"b\":null},{\"c\":113,\"t\":\"3cfade\",\"b\":null}]}",
	})
	// x := scrabble.NewGame("english", p, scrabble.NewDynamoValidator())
	// z := x.JSON()
	// fmt.Println(z)
}
