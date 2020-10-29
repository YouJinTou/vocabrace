package pooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/memcachier/gomemcache/memcache"
)

// Request encapsulates pool data.
type Request struct {
	ConnectionID string
	UserID       string
	PoolLimit    int
}

// JoinOrCreate adds a user to an existing pool
// (relative to their skill level), or creates a new one.
func (c Context) JoinOrCreate(r *Request) (*Pool, error) {
	var pool *Pool
	var err error
	bucket := c.getPoolBucket(&r.UserID)

	for i := 0; i < 30; i++ {
		if pool, err = c.mapConnectionToPool(bucket, r); err != nil {
			pool = c.newPool(bucket)

			c.updateBucket(bucket, pool.ID)
		} else {
			return pool, err
		}
	}

	return pool, err
}

func (c Context) getPoolBucket(userID *string) string {
	if userID == nil {
		return _Beginner
	}

	// Look up user's level
	return _Novice
}

func (c Context) getAvailablePoolID(bucket string) *string {
	c.minimizeRaceConditions()

	item, err := c.mc.Get(fmt.Sprintf("%s|currentAvailablePool", bucket))

	if item == nil || err != nil {
		return nil
	}

	poolID := string(item.Value)

	return &poolID
}

func (c Context) minimizeRaceConditions() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (c Context) getPool(poolID *string) (*Pool, *memcache.Item) {
	if poolID == nil {
		return nil, nil
	}

	item, err := c.mc.Get(*poolID)

	if err != nil {
		return nil, item
	}

	var connections []string

	json.Unmarshal(item.Value, &connections)

	return &Pool{
		ID:            *poolID,
		ConnectionIDs: connections,
	}, item
}

func (c Context) newPool(bucket string) *Pool {
	poolID := uuid.New().String()
	emptyListBytes, _ := json.Marshal([]string{})

	c.mc.Set(&memcache.Item{
		Key:   poolID,
		Value: emptyListBytes,
	})

	c.mc.Set(&memcache.Item{
		Key:   fmt.Sprintf("%s|currentAvailablePool", bucket),
		Value: []byte(poolID),
	})

	return &Pool{
		ID:            poolID,
		ConnectionIDs: []string{},
	}
}

func (c Context) mapConnectionToPool(bucket string, r *Request) (*Pool, error) {
	var pool *Pool
	var item *memcache.Item

	for {
		poolID := c.getAvailablePoolID(bucket)
		pool, item = c.getPool(poolID)

		if pool == nil || len(pool.ConnectionIDs) >= r.PoolLimit {
			return pool, errors.New("no suitable pool")
		}

		newConnections := append(pool.ConnectionIDs, r.ConnectionID)
		marshalled, _ := json.Marshal(newConnections)
		item.Value = marshalled
		casErr := c.mc.Cas(item)

		if casErr == nil {
			pool.ConnectionIDs = newConnections

			break
		}
	}

	setErr := c.mc.Set(&memcache.Item{
		Key:   r.ConnectionID,
		Value: []byte(pool.ID),
	})

	return pool, setErr
}

func (c Context) updateBucket(bucket, poolID string) error {
	key := fmt.Sprintf("%s|pools", bucket)
	_, err := c.mc.Get(key)

	if err != nil {
		empty, _ := json.Marshal([]string{})
		c.mc.Set(&memcache.Item{
			Key:   fmt.Sprintf("%s|pools", bucket),
			Value: empty,
		})
	}

	return c.mc.ListAppend(key, poolID)
}
