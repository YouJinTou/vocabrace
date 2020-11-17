package scrabble

import (
	"testing"
)

func TestCalculatePointsSumsTileValues(t *testing.T) {
	g := testScoreGame()
	w := word("test", 65, true, []int{}, []int{1, 2, 3, 4}, []int{})
	g.Board.SetCells(w.Cells)
	words := Extract(g.Board, w)
	scoreRunnable(g, w, words, 10, t)
}

func TestCalculatePointsCountsBlanksAsZeros(t *testing.T) {
	g := testScoreGame()
	w := word("test", 65, true, []int{}, []int{1, 2, 3, 4}, []int{0})
	g.Board.SetCells(w.Cells)
	words := Extract(g.Board, w)
	scoreRunnable(g, w, words, 9, t)
}

func TestCalculatePointsOriginDoubles(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	words := Extract(g.Board, w)
	scoreRunnable(g, w, words, 14, t)
}

func TestCalculatePointsPremiumAlreadyUsed(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)

	w1 := word("farm", 82, false, []int{2}, []int{4, 1, 1, 3}, []int{})
	g.Board.SetCells(w1.Cells)
	words := Extract(g.Board, w1)
	scoreRunnable(g, w1, words, 9, t)
}

func TestCalculatePointsCreateTwoNewWordsWithLetter(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	w1 := word("farm", 82, false, []int{2}, []int{4, 1, 1, 3}, []int{})
	g.Board.SetCells(w1.Cells)

	w3 := word("paste", 140, true, []int{}, []int{3, 1, 1, 1}, []int{})
	g.Board.SetCells(w3.Cells)
	words := Extract(g.Board, w3)
	scoreRunnable(g, w3, words, 25, t)
}

func TestCalculatePointsCreateThreeNewWordsWithPremiumLetter(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	w1 := word("farm", 82, false, []int{2}, []int{4, 1, 1, 3}, []int{})
	g.Board.SetCells(w1.Cells)
	w3 := word("paste", 140, true, []int{}, []int{3, 1, 1, 1}, []int{})
	g.Board.SetCells(w3.Cells)

	w4 := word("mob", 127, true, []int{0}, []int{3, 1, 3}, []int{})
	g.Board.SetCells(w4.Cells)
	words := Extract(g.Board, w4)
	scoreRunnable(g, w4, words, 16, t)
}

func TestCalculatePointsCreateThreeNewWordsWithPremiumWord(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	w1 := word("farm", 82, false, []int{2}, []int{4, 1, 1, 3}, []int{})
	g.Board.SetCells(w1.Cells)
	w3 := word("paste", 140, true, []int{}, []int{3, 1, 1, 1}, []int{})
	g.Board.SetCells(w3.Cells)
	w4 := word("mob", 127, true, []int{0}, []int{3, 1, 3}, []int{})
	g.Board.SetCells(w4.Cells)

	w5 := word("bit", 154, true, []int{}, []int{3, 1, 1}, []int{})
	g.Board.SetCells(w5.Cells)
	words := Extract(g.Board, w5)
	scoreRunnable(g, w5, words, 16, t)
}

func testScoreGame() *Game {
	return NewGame([]*Player{testPlayer(), testPlayer()})
}

func scoreRunnable(g *Game, w *Word, words []*Word, expected int, t *testing.T) {
	p := CalculatePoints(g, w, words)

	if p != expected {
		t.Errorf("expected %d, got %d", expected, p)
	}
}
