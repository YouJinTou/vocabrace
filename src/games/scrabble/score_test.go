package scrabble

import (
	"testing"
)

func TestCalculatePointsSumsTileValues(t *testing.T) {
	g := testScoreGame()
	w := word("test", 65, true, []int{}, []int{1, 2, 3, 4}, []int{})
	g.Board.SetCells(w.Cells)
	words := Extract(&g.Board, w)
	scoreRunnable(g, w, words, 10, t)
}

func TestCalculatePointsAwards50PointsIfAll7TilesUsed(t *testing.T) {
	g := testScoreGame()
	w := word("testing", 65, true, []int{}, []int{1, 1, 1, 1, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	words := Extract(&g.Board, w)
	scoreRunnable(g, w, words, (7*2)+50, t)
}

func TestCalculatePointsCountsBlanksAsZeros(t *testing.T) {
	g := testScoreGame()
	w := word("test", 65, true, []int{}, []int{1, 2, 3, 4}, []int{0})
	g.Board.SetCells(w.Cells)
	words := Extract(&g.Board, w)
	scoreRunnable(g, w, words, 9, t)
}

func TestCalculatePointsOriginDoubles(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)
	words := Extract(&g.Board, w)
	scoreRunnable(g, w, words, 14, t)
}

func TestCalculatePointsPremiumAlreadyUsed(t *testing.T) {
	g := testScoreGame()
	w := word("horn", 110, true, []int{}, []int{4, 1, 1, 1}, []int{})
	g.Board.SetCells(w.Cells)

	w1 := word("farm", 82, false, []int{2}, []int{4, 1, 1, 3}, []int{})
	g.Board.SetCells(w1.Cells)
	words := Extract(&g.Board, w1)
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
	words := Extract(&g.Board, w3)
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
	words := Extract(&g.Board, w4)
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
	words := Extract(&g.Board, w5)
	scoreRunnable(g, w5, words, 16, t)
}

func TestCalculatePoints_StartBingo(t *testing.T) {
	g := testScoreGame()
	w := testCreateWord(BoardOrigin, true, g.ToMove().Tiles.Value...)
	w.Cells[0].Tile.Value = 1
	w.Cells[1].Tile.Value = 4
	w.Cells[2].Tile.Value = 2
	w.Cells[3].Tile.Value = 2
	w.Cells[4].Tile.Value = 1
	w.Cells[5].Tile.Value = 1
	w.Cells[6].Tile.Value = 1
	g.Board.SetCells(w.Cells)
	words := Extract(&g.Board, w)
	scoreRunnable(g, w, words, 76, t)
}

func TestCalculatePoints_Bingo_DoubleWord_TripleWord_DoubleLetter(t *testing.T) {
	g := testScoreGame()
	w := word("farm", 117, false, []int{}, []int{}, []int{})
	g.Board.SetCells(w.Cells)

	w1 := testCreateWord(BoardOrigin, true, g.ToMove().Tiles.Value...)
	w1.Cells[0].Tile.Value = 1
	w1.Cells[1].Tile.Value = 4
	w1.Cells[2].Tile.Value = 2
	w1.Cells[3].Tile.Value = 2
	w1.Cells[4].Tile.Value = 3
	w1.Cells[5].Tile.Value = 1
	w1.Cells[5].Index = 118
	w1.Cells[6].Index = 119
	w1.Cells[6].Tile.Value = 2
	g.Board.SetCells(w1.Cells)

	tilesSum := 19
	multipliers := 2 * 3
	bingo := 50
	expected := (tilesSum * multipliers) + bingo
	words := Extract(&g.Board, w1)
	scoreRunnable(g, w1, words, expected, t)
}

func testScoreGame() *Game {
	return NewGame(English, []*Player{testPlayer(), testPlayer()}, v())
}

func scoreRunnable(g *Game, w *Word, words []*Word, expected int, t *testing.T) {
	p := CalculatePoints(w, words)

	if p != expected {
		t.Errorf("expected %d, got %d", expected, p)
	}
}
