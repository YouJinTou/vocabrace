package scrabble

import (
	"github.com/YouJinTou/vocabrace/tools"
)

// CalculatePoints calculates points given some words.
func CalculatePoints(primary *Word, words []*Word) int {
	disableMultipliers(primary, words)
	sum := sumTiles(words)
	sum = tryAward50PointsBonus(sum, primary)
	return sum
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

func tryAward50PointsBonus(sum int, primary *Word) int {
	if primary.Length() == 7 {
		return sum + 50
	}
	return sum
}
