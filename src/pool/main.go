package pool

import (
	"github.com/YouJinTou/vocabrace/memcached"
)

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
