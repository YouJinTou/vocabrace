package scrabble

type language func() []Tile

// English produces a full set of tiles for the Latin alphabet.
func English() []Tile {
	return []Tile{
		*BlankTile(),
		Tile{"A", 1},
		Tile{"B", 3},
		Tile{"C", 3},
		Tile{"D", 2},
		Tile{"E", 1},
		Tile{"F", 4},
		Tile{"G", 2},
		Tile{"H", 4},
		Tile{"I", 1},
		Tile{"J", 8},
		Tile{"K", 5},
		Tile{"L", 1},
		Tile{"M", 3},
		Tile{"N", 1},
		Tile{"O", 1},
		Tile{"P", 3},
		Tile{"R", 1},
		Tile{"S", 1},
		Tile{"T", 1},
		Tile{"U", 1},
		Tile{"V", 4},
		Tile{"W", 4},
		Tile{"X", 1},
		Tile{"Y", 8},
		Tile{"Z", 10},
	}
}
