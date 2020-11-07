package memcachedpooling

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/YouJinTou/vocabrace/pooling"
	"github.com/YouJinTou/vocabrace/tools"

	"github.com/YouJinTou/vocabrace/memcached"
)

// MemcachedProvider handles adding and removing connections from pools.
type MemcachedProvider struct {
	mc *memcached.Client
}

// NewMemcachedProvider creates a new pooling provider using Memcached as a backend.
func NewMemcachedProvider(host, username, password string) pooling.Provider {
	return MemcachedProvider{
		mc: memcached.New(host, username, password),
	}
}

// GetPool gets a pool by ID.
func (c MemcachedProvider) GetPool(i *pooling.GetPoolInput) (*pooling.Pool, error) {
	p, _ := c.getPool(&i.PoolID)

	if p == nil {
		return nil, errors.New("pool not found")
	}

	return p, nil
}

// GetPeers maps a connectionID to a pool and returns all peer connections.
func (c MemcachedProvider) GetPeers(i *pooling.GetPeersInput) ([]string, error) {
	it, err := c.mc.Get(i.ConnectionID)

	if err != nil {
		return nil, fmt.Errorf("not found for connection %s", i.ConnectionID)
	}

	poolID := string(it.Value)
	getIt, getErr := c.mc.Get(poolID)

	if getErr != nil {
		return nil, fmt.Errorf("not found for pool %s", poolID)
	}

	connectionIDs := []string{}

	json.Unmarshal(getIt.Value, &connectionIDs)

	connectionIDs = tools.SliceRemoveString(connectionIDs, i.ConnectionID)

	return connectionIDs, nil
}

// List lists all pools
func (c MemcachedProvider) List() {
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
