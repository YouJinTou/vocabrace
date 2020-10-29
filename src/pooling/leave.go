package pooling

// Leave removes a connection from a given pool.
func (c Context) Leave(connectionID string) (*Pool, error) {
	item, err := c.mc.Get(connectionID)

	if err != nil {
		return nil, err
	}

	poolID := string(item.Value)
	_, removeErr := c.mc.ListRemove(poolID, connectionID)

	if removeErr != nil {
		return nil, removeErr
	}

	deleteErr := c.mc.Delete(connectionID)
	pool, _ := c.getPool(&poolID)

	return pool, deleteErr
}
