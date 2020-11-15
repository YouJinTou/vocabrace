package scrabble

type language func() Tiles

// English produces a full set of tiles for the Latin alphabet.
func English() Tiles {
	return *NewTiles(
		&Tile{"A", 1, tileID()},
		&Tile{"B", 3, tileID()},
		&Tile{"C", 3, tileID()},
		&Tile{"D", 2, tileID()},
		&Tile{"E", 1, tileID()},
		&Tile{"F", 4, tileID()},
		&Tile{"G", 2, tileID()},
		&Tile{"H", 4, tileID()},
		&Tile{"I", 1, tileID()},
		&Tile{"J", 8, tileID()},
		&Tile{"K", 5, tileID()},
		&Tile{"L", 1, tileID()},
		&Tile{"M", 3, tileID()},
		&Tile{"N", 1, tileID()},
		&Tile{"O", 1, tileID()},
		&Tile{"P", 3, tileID()},
		&Tile{"Q", 10, tileID()},
		&Tile{"R", 1, tileID()},
		&Tile{"S", 1, tileID()},
		&Tile{"T", 1, tileID()},
		&Tile{"U", 1, tileID()},
		&Tile{"V", 4, tileID()},
		&Tile{"W", 4, tileID()},
		&Tile{"X", 8, tileID()},
		&Tile{"Y", 4, tileID()},
		&Tile{"Z", 10, tileID()},
	)
}
