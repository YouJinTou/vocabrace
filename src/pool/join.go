package pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
)

// Request encapsulates pool data.
type Request struct {
	ConnectionID string
	UserID       string
	PoolLimit    int
}

type pool struct {
	ID          string
	Connections []string
}

func (p pool) isFull(limit int) bool {
	return len(p.Connections) >= limit
}

// JoinOrCreate adds a user to an existing pool
// (relative to their skill level), or creates a new one.
func (p Pool) JoinOrCreate(r *Request) {
	bucket := p.getPoolBucket(&r.UserID)

	for {
		var poolID = p.getAvailablePool(bucket)
		var pool = p.getPool(poolID)

		if pool == nil {
			pool = p.newPool(bucket)

			p.updateBucket(bucket, pool.ID)
		}

		if err := p.mapConnectionToPool(r.ConnectionID, pool, r.PoolLimit); err != nil {
			pool = p.newPool(bucket)

			p.updateBucket(bucket, pool.ID)
		} else {
			break
		}
	}
}

func (p Pool) getPoolBucket(userID *string) string {
	if userID == nil {
		return _Beginner
	}

	// Look up user's level
	return _Novice
}

func (p Pool) getAvailablePool(bucket string) *string {
	p.minimizeRaceConditions()

	item, err := p.c.Get(fmt.Sprintf("%s|currentAvailablePool", bucket))

	if item == nil || err != nil {
		return nil
	}

	poolID := string(item.Value)

	return &poolID
}

func (p Pool) minimizeRaceConditions() {
	time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
}

func (p Pool) getPool(poolID *string) *pool {
	if poolID == nil {
		return nil
	}

	item, err := p.c.Get(*poolID)

	if err != nil {
		return nil
	}

	var connections []string

	json.Unmarshal(item.Value, &connections)

	return &pool{
		ID:          *poolID,
		Connections: connections,
	}
}

func (p Pool) newPool(bucket string) *pool {
	poolID := uuid.New().String()
	emptyListBytes, _ := json.Marshal([]string{})

	p.c.Set(&memcache.Item{
		Key:   poolID,
		Value: emptyListBytes,
	})

	p.c.Set(&memcache.Item{
		Key:   fmt.Sprintf("%s|currentAvailablePool", bucket),
		Value: []byte(poolID),
	})

	return &pool{
		ID:          poolID,
		Connections: []string{},
	}
}

func (p Pool) mapConnectionToPool(connectionID string, pool *pool, poolLimit int) error {
	for i := 0; i < 10; i++ {
		item, getErr := p.c.Get(pool.ID)

		if getErr != nil {
			continue
		}

		var oldItems []string

		json.Unmarshal(item.Value, &oldItems)

		if len(oldItems) >= poolLimit {
			return errors.New("pool is full")
		}

		newItems := append(oldItems, connectionID)
		newItemsMarshalled, _ := json.Marshal(newItems)
		item.Value = newItemsMarshalled
		casErr := p.c.Cas(item)

		if casErr == nil {
			break
		}
	}

	setErr := p.c.Set(&memcache.Item{
		Key:   connectionID,
		Value: []byte(pool.ID),
	})

	return setErr
}

func (p Pool) updateBucket(bucket, poolID string) error {
	key := fmt.Sprintf("%s|pools", bucket)
	_, err := p.c.Get(key)

	if err != nil {
		empty, _ := json.Marshal([]string{})
		p.c.Set(&memcache.Item{
			Key:   fmt.Sprintf("%s|pools", bucket),
			Value: empty,
		})
	}

	return p.c.ListAppend(key, poolID)
}
