package scrabble

// CalculatePoints calculates points given some words.
func CalculatePoints(g *Game, primary *Word, words []*Word) int {
	return calculateTilesSum(words)
}

func calculateTilesSum(words []*Word) int {
	sum := 0
	for _, w := range words {
		sum += w.Value()
	}
	return sum
}
