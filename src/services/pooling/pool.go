package pooling

// Pool carries pool data.
type Pool struct {
	ID            string
	ConnectionIDs []string
	Bucket        string
	Limit         int
}

// Beginner is the beginner bucket.
const Beginner = "beginner"

// Novice is the novice bucket.
const Novice = "novice"

// LowerIntermediate is the lower_intermediate bucket.
const LowerIntermediate = "lower_intermediate"

// Intermediate is the intermediate bucket.
const Intermediate = "intermediate"

// UpperIntermediate is the upper_intermediate bucket.
const UpperIntermediate = "upper_intermediate"

// Advanced is the advanced bucket.
const Advanced = "advanced"

// Expert is the expert bucket.
const Expert = "expert"

// Godlike is the godlike bucket.
const Godlike = "godlike"
