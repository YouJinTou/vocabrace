package pooling

import (
	"encoding/json"
	"fmt"

	"github.com/YouJinTou/vocabrace/core"

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

// Context handles adding and removing connections from pools.
type Context struct {
	mc *memcached.Client
}

// Pool carries pool data.
type Pool struct {
	ID            string
	ConnectionIDs []string
}

// NewMemcachedContext creates a new context using Memcached as a backend.
func NewMemcachedContext(host, username, password string) Context {
	return Context{
		mc: memcached.New(host, username, password),
	}
}

// GetPeers maps a connectionID to a pool and returns all peer connections.
func (c Context) GetPeers(connectionID string) ([]string, error) {
	it, err := c.mc.Get(connectionID)

	if err != nil {
		return nil, fmt.Errorf("not found for connection %s", connectionID)
	}

	poolID := string(it.Value)
	getIt, getErr := c.mc.Get(poolID)

	if getErr != nil {
		return nil, fmt.Errorf("not found for pool %s", poolID)
	}

	connectionIDs := []string{}

	json.Unmarshal(getIt.Value, &connectionIDs)

	connectionIDs = core.SliceRemoveString(connectionIDs, connectionID)

	return connectionIDs, nil
}

// List lists all pools
func (c Context) List() {
	item, err := c.mc.Get("novice|pools")

	if err != nil {
		return
	}

	pools := []string{}
	json.Unmarshal(item.Value, &pools)

	fmt.Println(fmt.Sprintf("Total pools: %d", len(pools)))

	for _, curr := range pools {
		item, _ := c.mc.Get(curr)

		connections := []string{}
		json.Unmarshal(item.Value, &connections)

		fmt.Println(fmt.Sprintf("%s: %d", curr, len(connections)))
	}
}
