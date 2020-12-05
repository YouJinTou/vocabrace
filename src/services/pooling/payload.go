package pooling

// PoolerPayload encapsulates data needed to push a message to the pooler.
type PoolerPayload struct {
	Domain       string
	ConnectionID string
	Bucket       string
	Game         string
	Players      int
	Language     string
	UserID       string
	Username     string
}
