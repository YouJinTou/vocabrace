package scrabble

import (
	"github.com/YouJinTou/vocabrace/tools"
)

// CalculatePoints calculates points given some words.
func CalculatePoints(g *Game, primary *Word, words []*Word) int {
	disableMultipliers(primary, words)
	return sumTiles(words)
}

func disableMultipliers(primary *Word, words []*Word) {
	primaryIndices := primary.Indices()
	for _, w := range words {
		for _, wc := range w.Cells {
			if !tools.ContainsInt(primaryIndices, wc.Index) {
				wc.enableMultiplier = false
			}
		}
	}
}

func sumTiles(words []*Word) int {
	sum := 0
	for _, w := range words {
		sum += w.Value()
	}
	return sum
}
