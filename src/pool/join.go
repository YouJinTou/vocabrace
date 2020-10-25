package pool

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

type pool struct {
	ID          string
	Connections []string
	Item        *memcache.Item
}

// JoinOrCreate adds a user to an existing pool
// (relative to their skill level), or creates a new one.
func (p Pool) JoinOrCreate(r *Request) error {
	var err error
	bucket := p.getPoolBucket(&r.UserID)

	for i := 0; i < 30; i++ {
		if err = p.mapConnectionToPool(bucket, r); err != nil {
			pool := p.newPool(bucket)

			p.updateBucket(bucket, pool.ID)
		} else {
			break
		}
	}

	return err
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
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
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
		Item:        item,
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

func (p Pool) mapConnectionToPool(bucket string, r *Request) error {
	var poolID *string

	for {
		poolID = p.getAvailablePool(bucket)
		pool := p.getPool(poolID)

		if pool == nil || len(pool.Connections) >= r.PoolLimit {
			return errors.New("no suitable pool")
		}

		newConnections := append(pool.Connections, r.ConnectionID)
		marshalled, _ := json.Marshal(newConnections)
		pool.Item.Value = marshalled
		casErr := p.c.Cas(pool.Item)

		if casErr == nil {
			break
		}
	}

	setErr := p.c.Set(&memcache.Item{
		Key:   r.ConnectionID,
		Value: []byte(*poolID),
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
