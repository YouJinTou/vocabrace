package pool

import (
	"encoding/json"
	"fmt"

	"github.com/YouJinTou/vocabrace/memcached"
)

const _Beginner = "beginner"
const _Novice = "novice"
const _LowerIntermediate = "lower_intermediate"
const _Intermediate = "intermediate"
const _UpperIntermediate = "intermediate"
const _Advanced = "advanced"
const _Expert = "expert"
const _Godlike = "godlike"

// Pool handles adding and removing connections from pools.
type Pool struct {
	c *memcached.Client
}

// NewMemcached creates a new pool using Memcached as a backend.
func NewMemcached(host, username, password string) Pool {
	return Pool{
		c: memcached.New(host, username, password),
	}
}

// GetPeers maps a connectionID to a pool and returns all peer connections.
func (p Pool) GetPeers(connectionID string) ([]string, error) {
	it, err := p.c.Get(connectionID)

	if err != nil {
		return nil, fmt.Errorf("not found for connection %s", connectionID)
	}

	poolID := string(it.Value)
	getIt, getErr := p.c.Get(poolID)

	if getErr != nil {
		return nil, fmt.Errorf("not found for pool %s", poolID)
	}

	connectionIDs := []string{}
	json.Unmarshal(getIt.Value, &connectionIDs)

	for i, curr := range connectionIDs {
		if curr == connectionID {
			connectionIDs = append(connectionIDs[:i], connectionIDs[i+1:]...)
			break
		}
	}

	return connectionIDs, nil
}

// List lists all pools
func (p Pool) List() {
	item, err := p.c.Get("novice|pools")

	if err != nil {
		return
	}

	pools := []string{}
	json.Unmarshal(item.Value, &pools)

	fmt.Println(fmt.Sprintf("Total pools: %d", len(pools)))

	for _, curr := range pools {
		item, _ := p.c.Get(curr)

		connections := []string{}
		json.Unmarshal(item.Value, &connections)

		fmt.Println(fmt.Sprintf("%s: %d", curr, len(connections)))
	}
}
