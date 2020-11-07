package memcachedpooling

import "github.com/YouJinTou/vocabrace/pooling"

// Leave removes a connection from a given pool.
func (c MemcachedProvider) Leave(i *pooling.LeaveInput) (*pooling.Pool, error) {
	item, err := c.mc.Get(i.ConnectionID)

	if err != nil {
		return nil, err
	}

	poolID := string(item.Value)
	_, removeErr := c.mc.ListRemove(poolID, i.ConnectionID)

	if removeErr != nil {
		return nil, removeErr
	}

	deleteErr := c.mc.Delete(i.ConnectionID)
	pool, _ := c.getPool(&poolID)

	return pool, deleteErr
}
