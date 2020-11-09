package main

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func main() {
	p := []scrabble.Player{
		scrabble.Player{
			ID:     "1",
			Name:   "Name",
			Points: 1,
		},
	}
	x := scrabble.NewGame(p)
	z := x.JSON()
	fmt.Println(z)
}
