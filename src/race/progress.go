package race

var progress = make(map[string]int)

// GetProgress gets the current progress of all players
func GetProgress() map[string]int {
	return progress
}

// UpdateProgress updates a player's progress for each solved question
func UpdateProgress(userID string) {
	if _, ok := progress[userID]; ok {
		progress[userID]++
	} else {
		progress[userID] = 1
	}
}
