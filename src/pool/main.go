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

// New creates a new pool.
func New() Pool {
	return Pool{
		c: memcached.New("localhost:11211"),
	}
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
