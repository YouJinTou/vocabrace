package pool

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/google/uuid"
)

const _Beginner = "beginner"
const _Novice = "novice"
const _LowerIntermediate = "lower_intermediate"
const _Intermediate = "intermediate"
const _UpperIntermediate = "intermediate"
const _Advanced = "advanced"
const _Expert = "expert"
const _Godlike = "godlike"

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
	var poolID = p.getAvailablePool(bucket)
	var pool = p.getPool(poolID)

	if poolID == nil || pool.isFull(r.PoolLimit) {
		pool = p.newPool(bucket)

		p.updateBucket(bucket, pool.ID)
	}

	p.mapConnectionToPool(r.ConnectionID, pool)
}

func (p Pool) getPoolBucket(userID *string) string {
	if userID == nil {
		return _Beginner
	}

	// Look up user's level
	return _Novice
}

func (p Pool) getAvailablePool(bucket string) *string {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

	item, err := p.c.Get(fmt.Sprintf("%s|currentAvailablePool", bucket))

	if item == nil || err != nil {
		return nil
	}

	poolID := string(item.Value)

	return &poolID
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

func (p Pool) mapConnectionToPool(connectionID string, pool *pool) error {
	if err := p.c.Set(&memcache.Item{
		Key:   connectionID,
		Value: []byte(pool.ID),
	}); err != nil {
		return err
	}

	return p.c.ListAppend(pool.ID, connectionID)
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
