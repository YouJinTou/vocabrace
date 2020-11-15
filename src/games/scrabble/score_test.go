package scrabble

import (
	"testing"

	"github.com/YouJinTou/vocabrace/tools"
)

func TestCalculatePointsSumsTileValues(t *testing.T) {
	g := testScoreGame()
	w := scoreWord("test", []int{1, 2, 3, 4}, []int{})
	w2 := scoreWord("tests", []int{1, 2, 3, 4, 5}, []int{})
	expected := 25
	p := CalculatePoints(g, w, []*Word{w, w2})

	if p != expected {
		t.Errorf("expected %d, got %d", expected, p)
	}
}

func TestCalculatePointsCountsBlanksAsZeros(t *testing.T) {
	g := testScoreGame()
	wWithBlank := scoreWord("test", []int{1, 2, 3, 4}, []int{2})
	w2 := scoreWord("tests", []int{1, 2, 3, 4, 5}, []int{})
	expected := 22
	p := CalculatePoints(g, wWithBlank, []*Word{wWithBlank, w2})

	if p != expected {
		t.Errorf("expected %d, got %d", expected, p)
	}
}

func testScoreGame() *Game {
	return NewGame([]*Player{testPlayer(), testPlayer()})
}

func scoreWord(w string, values, blanks []int) *Word {
	if len(w) != len(values) {
		panic("invalid score word")
	}

	cells := []*Cell{}
	for i, ch := range w {
		var t *Tile
		if tools.ContainsInt(blanks, i) {
			t = BlankTile()
		} else {
			t = NewTile(string(ch), values[i])
		}
		c := NewCell(t, i)
		cells = append(cells, c)
	}
	return NewWord(cells)
}
