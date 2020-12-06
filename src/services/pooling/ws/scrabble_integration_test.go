package ws

import (
	"os"
	"testing"

	"github.com/YouJinTou/vocabrace/games/scrabble"
)

func TestDoublesPointsOnFirstMove(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("STAGE", "dev")
	os.Setenv("AWS_PROFILE", "vocabrace")

	pid := "testing_pid"
	users := []*User{&User{}, &User{}}
	sws := scrabblews{}
	players := sws.loadPlayers(users)
	g := scrabble.NewGame("bulgarian", players, scrabble.NewDynamoValidator())
	w := scrabble.NewWordFromString(
		"ТИ", []int{2, 3}, []int{scrabble.BoardOrigin, scrabble.BoardOrigin + 1})
	g.ToMove().Tiles.RemoveAt(0)
	g.ToMove().Tiles.RemoveAt(0)
	g.ToMove().Tiles.Append(&w.Cells[0].Tile, &w.Cells[1].Tile)

	if _, err := g.Place(w); err != nil {
		t.Errorf(err.Error())
	}

	if err := saveState(pid, g); err != nil {
		t.Errorf(err.Error())
	}

	game := &scrabble.Game{}
	loadState(pid, game)

	if game.GetPlayerByID(g.GetLastMovedID()).Points != 10 {
		t.Errorf("invalid points")
	}
}
