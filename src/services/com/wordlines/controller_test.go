package controller

import (
	"testing"

	"github.com/YouJinTou/vocabrace/games/wordlines"
)

type canValidateMock struct{}

func (m canValidateMock) ValidatePlace(g wordlines.Game, w *wordlines.Word) error {
	return nil
}

func TestPlace(t *testing.T) {
	s := Controller{}
	players := []*wordlines.Player{&wordlines.Player{}}
	g := wordlines.NewClassicGame("english", players, canValidateMock{})
	tu := &turn{
		IsPlace: true,
		Word: []*cell{
			&cell{
				CellIndex: 112,
				TileID:    g.ToMove().Tiles.GetAt(0).ID,
			},
			&cell{
				CellIndex: 113,
				TileID:    g.ToMove().Tiles.GetAt(1).ID,
			},
		},
	}
	result := s.place(tu, g)

	if result.err != nil {
		t.Errorf("did not expect error")
	}
}
