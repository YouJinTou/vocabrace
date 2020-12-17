package wordlines

type layout interface {
	DoubleLetterIndices() []int
	TripleLetterIndices() []int
	TripleWordIndices() []int
	DoubleWordIndices() []int
}

type classic struct{}
type spiral struct{}

func (c classic) DoubleLetterIndices() []int {
	return []int{3, 11, 36, 38, 45, 52, 59, 92, 96, 98, 102, 108, 116, 122, 126, 128, 132, 165, 172, 179, 186, 188, 213, 221}
}
func (c classic) TripleLetterIndices() []int {
	return []int{20, 24, 76, 80, 84, 88, 136, 140, 144, 148, 200, 204}
}
func (c classic) TripleWordIndices() []int {
	return []int{0, 7, 14, 105, 119, 210, 217, 224}
}
func (c classic) DoubleWordIndices() []int {
	return []int{16, 32, 48, 64, 112, 160, 176, 192, 208, 28, 42, 56, 70, 154, 168, 182, 196}
}

func (s spiral) DoubleLetterIndices() []int {
	return []int{50, 36, 86, 102, 118, 174, 188, 138, 122, 106, 81, 45, 61, 125, 141, 199, 143, 179, 163, 99, 83, 25, 213, 11}
}
func (s spiral) TripleLetterIndices() []int {
	return []int{4, 8, 39, 74, 134, 129, 220, 216, 185, 150, 90, 95}
}
func (s spiral) TripleWordIndices() []int {
	return []int{0, 14, 224, 210, 22, 202, 77, 147}
}
func (s spiral) DoubleWordIndices() []int {
	return []int{112, 16, 32, 48, 64, 28, 42, 56, 70, 208, 192, 176, 160, 196, 182, 168, 154}
}
