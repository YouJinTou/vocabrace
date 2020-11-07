package memcachedpooling

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/YouJinTou/vocabrace/pooling"

	"github.com/google/uuid"
	"github.com/memcachier/gomemcache/memcache"
)

// JoinOrCreate adds a user to an existing pool
// (relative to their skill level), or creates a new one.
func (c MemcachedProvider) JoinOrCreate(r *pooling.Request) (*pooling.Pool, error) {
	var pool *pooling.Pool
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

// GetPool gets a pool by ID.
func (c MemcachedProvider) GetPool(poolID string, r *pooling.Request) (*pooling.Pool, error) {
	p, _ := c.getPool(&poolID)

	if p == nil {
		return nil, errors.New("pool not found")
	}

	return p, nil
}

func (c MemcachedProvider) getPoolBucket(userID *string) string {
	if userID == nil {
		return pooling.Beginner
	}

	// Look up user's level
	return pooling.Novice
}

func (c MemcachedProvider) getAvailablePoolID(bucket string) *string {
	c.minimizeRaceConditions()

	item, err := c.mc.Get(fmt.Sprintf("%s|currentAvailablePool", bucket))

	if item == nil || err != nil {
		return nil
	}

	poolID := string(item.Value)

	return &poolID
}

func (c MemcachedProvider) minimizeRaceConditions() {
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
}

func (c MemcachedProvider) getPool(poolID *string) (*pooling.Pool, *memcache.Item) {
	if poolID == nil {
		return nil, nil
	}

	item, err := c.mc.Get(*poolID)

	if err != nil {
		return nil, item
	}

	var connections []string

	json.Unmarshal(item.Value, &connections)

	return &pooling.Pool{
		ID:            *poolID,
		ConnectionIDs: connections,
	}, item
}

func (c MemcachedProvider) newPool(bucket string) *pooling.Pool {
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

	return &pooling.Pool{
		ID:            poolID,
		ConnectionIDs: []string{},
	}
}

func (c MemcachedProvider) mapConnectionToPool(bucket string, r *pooling.Request) (*pooling.Pool, error) {
	var pool *pooling.Pool
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

func (c MemcachedProvider) updateBucket(bucket, poolID string) error {
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
