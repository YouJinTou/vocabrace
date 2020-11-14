package scrabble

import "github.com/google/uuid"

type language func() []Tile

// English produces a full set of tiles for the Latin alphabet.
func English() []Tile {
	return []Tile{
		Tile{"A", 1, uuid.New().String()},
		Tile{"B", 3, uuid.New().String()},
		Tile{"C", 3, uuid.New().String()},
		Tile{"D", 2, uuid.New().String()},
		Tile{"E", 1, uuid.New().String()},
		Tile{"F", 4, uuid.New().String()},
		Tile{"G", 2, uuid.New().String()},
		Tile{"H", 4, uuid.New().String()},
		Tile{"I", 1, uuid.New().String()},
		Tile{"J", 8, uuid.New().String()},
		Tile{"K", 5, uuid.New().String()},
		Tile{"L", 1, uuid.New().String()},
		Tile{"M", 3, uuid.New().String()},
		Tile{"N", 1, uuid.New().String()},
		Tile{"O", 1, uuid.New().String()},
		Tile{"P", 3, uuid.New().String()},
		Tile{"R", 1, uuid.New().String()},
		Tile{"S", 1, uuid.New().String()},
		Tile{"T", 1, uuid.New().String()},
		Tile{"U", 1, uuid.New().String()},
		Tile{"V", 4, uuid.New().String()},
		Tile{"W", 4, uuid.New().String()},
		Tile{"X", 1, uuid.New().String()},
		Tile{"Y", 8, uuid.New().String()},
		Tile{"Z", 10, uuid.New().String()},
	}
}
