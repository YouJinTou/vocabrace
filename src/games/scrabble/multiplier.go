package scrabble

type multiplier interface {
	Apply() int
}

type nilMultiplier struct{ t *Tile }

func (m *nilMultiplier) Apply() int { return m.t.Value }

type doubleLetterMultiplier struct{ t *Tile }

func (m *doubleLetterMultiplier) Apply() int { return m.t.Value * 2 }

type tripleLetterMultiplier struct{ t *Tile }

func (m *tripleLetterMultiplier) Apply() int { return m.t.Value * 3 }

type doubleWordMultiplier struct{ w *Word }

func (m *doubleWordMultiplier) Apply() int { return m.w.Value() * 2 }

type tripleWordMultiplier struct{ w *Word }

func (m *tripleWordMultiplier) Apply() int { return m.w.Value() * 3 }
