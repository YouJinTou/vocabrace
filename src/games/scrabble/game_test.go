package scrabble

import "testing"

func TestExchange(t *testing.T) {
	players := []*Player{testGetPlayer(), testGetPlayer()}
	g := NewGame(players)
	toExchange := []string{g.ToMove.Tiles[0].Letter, g.ToMove.Tiles[1].Letter}
	_, err := g.Exchange(toExchange)

	if err != nil {
		t.Errorf(err.Error())
	}
}
