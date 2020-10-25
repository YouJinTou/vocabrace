package pool

// Leave removes a connection from a given pool.
func (p Pool) Leave(connectionID string) error {
	item, err := p.c.Get(connectionID)

	if err != nil {
		return err
	}

	poolID := string(item.Value)
	removeErr := p.c.ListRemove(poolID, connectionID)

	if removeErr != nil {
		return removeErr
	}

	deleteErr := p.c.Delete(connectionID)

	return deleteErr
}
