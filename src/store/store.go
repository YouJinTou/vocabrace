package store

// Question wraps question data.
type Question struct {
	Text             string
	Answers          []string
	CorrectAnswerIdx int
}

// GetQuestions gets questions
func GetQuestions() []Question {
	return []Question{
		Question{
			Text:             "cat",
			Answers:          []string{"a four-legged critter", "the pope", "a small leafless plant"},
			CorrectAnswerIdx: 0,
		},
		Question{
			Text:             "dog",
			Answers:          []string{"a four-legged critter", "the pope", "a small leafless plant"},
			CorrectAnswerIdx: 0,
		},
		Question{
			Text:             "chivalry",
			Answers:          []string{"a four-legged critter", "the pope", "a small leafless plant"},
			CorrectAnswerIdx: 0,
		},
		Question{
			Text:             "newtonian",
			Answers:          []string{"a four-legged critter", "the pope", "a small leafless plant"},
			CorrectAnswerIdx: 0,
		},
		Question{
			Text:             "grass",
			Answers:          []string{"a four-legged critter", "the pope", "a small leafless plant"},
			CorrectAnswerIdx: 0,
		},
	}
}
