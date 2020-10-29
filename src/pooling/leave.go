package pooling

// Leave removes a connection from a given pool.
func (c Context) Leave(connectionID string) error {
	item, err := c.mc.Get(connectionID)

	if err != nil {
		return err
	}

	poolID := string(item.Value)
	removeErr := c.mc.ListRemove(poolID, connectionID)

	if removeErr != nil {
		return removeErr
	}

	deleteErr := c.mc.Delete(connectionID)

	return deleteErr
}
