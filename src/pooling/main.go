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

// JoinOrCreateInput is the input to JoinOrCreate.
type JoinOrCreateInput struct {
	ConnectionID string
	Bucket       string
	PoolLimit    int
}

// LeaveInput is the input to Leave.
type LeaveInput struct {
	ConnectionID string
	Bucket       string
}

// GetPoolInput is the input to GetPool.
type GetPoolInput struct {
	PoolID string
	Bucket string
}

// GetPeersInput is the input to GetPeers.
type GetPeersInput struct {
	ConnectionID string
	Bucket       string
}

// Provider abstracts a pooling provider.
type Provider interface {
	JoinOrCreate(i *JoinOrCreateInput) (*Pool, error)
	Leave(i *LeaveInput) (*Pool, error)
	GetPool(i *GetPoolInput) (*Pool, error)
	GetPeers(i *GetPeersInput) ([]string, error)
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
