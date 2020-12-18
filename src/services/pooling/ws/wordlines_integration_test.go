package ws

import (
	"os"
	"testing"

	"github.com/YouJinTou/vocabrace/games/wordlines"
)

func TestDoublesPointsOnFirstMove(t *testing.T) {
	os.Setenv("AWS_REGION", "eu-central-1")
	os.Setenv("STAGE", "dev")
	os.Setenv("AWS_PROFILE", "vocabrace")

	pid := "testing_pid"
	sws := wordlinesws{}
	cons := []*Connection{&Connection{}, &Connection{}}
	players := sws.loadPlayers(NewConnections(cons))
	g := wordlines.NewClassicGame("bulgarian", players, wordlines.NewDynamoValidator())
	w := wordlines.NewWordFromString(
		"ТИ", []int{2, 3}, []int{wordlines.BoardOrigin, wordlines.BoardOrigin + 1})
	g.ToMove().Tiles.RemoveAt(0)
	g.ToMove().Tiles.RemoveAt(0)
	g.ToMove().Tiles.Append(&w.Cells[0].Tile, &w.Cells[1].Tile)

	if _, err := g.Place(w); err != nil {
		t.Errorf(err.Error())
	}

	if err := saveState(&saveStateInput{
		PoolID:        pid,
		ConnectionIDs: []string{"123", "456"},
		V:             g,
	}); err != nil {
		t.Errorf(err.Error())
	}

	game := &wordlines.Game{}
	loadState(pid, game)

	if game.GetPlayerByID(g.GetLastMovedID()).Points != 10 {
		t.Errorf("invalid points")
	}
}
