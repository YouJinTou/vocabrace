package main

import (
	"fmt"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func main() {
	g := scrabble.NewGame()
	gameStr := g.JSON()
	fmt.Println(gameStr)
}
