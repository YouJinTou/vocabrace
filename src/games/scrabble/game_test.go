package scrabble

import (
	"strconv"
	"testing"
)

func TestExchange(t *testing.T) {
	players := []*Player{testPlayer(), testPlayer()}
	g := NewGame(players)
	toExchange := []string{g.ToMove().Tiles[0].Letter, g.ToMove().Tiles[1].Letter}
	_, err := g.Exchange(toExchange)

	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestOrderPlayers(t *testing.T) {
	for total := 1; total < 50; total++ {
		players := []*Player{}
		for i := 0; i < total; i++ {
			players = append(players, testPlayerArgs(strconv.Itoa(i)))
		}

		for x := 0; x < 50; x++ {
			g := NewGame(players)
			io, _ := strconv.Atoi(g.Order[0])
			expected := getExpectedOrder(io, total)

			if len(expected) != len(g.Order) {
				t.Errorf("Lengths mismatch.")
			}

			for i := 0; i < len(g.Order); i++ {
				if expected[i] != g.Order[i] {
					t.Errorf("Invalid order.")
				}
			}
		}
	}
}

func TestSetNext(t *testing.T) {
	p1 := testPlayerArgs("1")
	p2 := testPlayerArgs("2")
	p3 := testPlayerArgs("3")
	players := []*Player{p1, p2, p3}

	for j := 0; j < 10; j++ {
		g := NewGame(players)

		for i := 0; i < 50; i++ {
			toMoveID := g.ToMoveID
			g.setNext()
			failed := false

			if toMoveID == p1.ID {
				if g.ToMoveID != p2.ID {
					failed = true
				}
			} else if toMoveID == p2.ID {
				if g.ToMoveID != p3.ID {
					failed = true
				}
			} else {
				if g.ToMoveID != p1.ID {
					failed = true
				}
			}

			if failed {
				t.Errorf("could not set next player")
			}
		}
	}
}

func getExpectedOrder(idx, total int) []string {
	result := []string{strconv.Itoa(idx)}
	for i := idx + 1; i < total; i++ {
		result = append(result, strconv.Itoa(i))
	}
	for i := 0; i < idx; i++ {
		result = append(result, strconv.Itoa(i))
	}
	return result
}
