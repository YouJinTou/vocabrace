package ws

// OnStart executes communication logic at the start of a game.
func OnStart(data *ReceiverData) {
	switch data.Game {
	case "scrabble":
		scrabbleOnStart(data)
	default:
		panic("invalid game")
	}
}
