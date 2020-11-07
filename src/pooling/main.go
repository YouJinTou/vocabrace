package pooling

// Pool carries pool data.
type Pool struct {
	ID            string
	ConnectionIDs []string
	Bucket        string
	Limit         int
}

// GetPeers gets a connection's pool peers.
func (p *Pool) GetPeers(connectionID string) []string {
	peers := []string{}

	for _, cid := range p.ConnectionIDs {
		if cid != connectionID {
			peers = append(peers, cid)
		}
	}

	return peers
}

// Request encapsulates pool request data.
type Request struct {
	ConnectionID string
	UserID       string
	Bucket       string
	PoolLimit    int
	Stage        string
}

// Provider abstracts a pooling provider.
type Provider interface {
	JoinOrCreate(r *Request) (*Pool, error)
	Leave(r *Request) (*Pool, error)
	GetPool(poolID string, r *Request) (*Pool, error)
	GetPeers(request *Request) ([]string, error)
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
