package memcached

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/memcachier/gomemcache/memcache"
)

// Client wraps the bradfitz's memcache.Client.
type Client struct {
	bfc *memcache.Client
}

// New creates a client with credentials
func New(server, username, password string) *Client {
	client := memcache.New(server)
	client.SetAuth(username, []byte(password))

	return &Client{
		bfc: client,
	}
}

// Get sets a key. Will retry 3 times.
func (c Client) Get(key string) (*memcache.Item, error) {
	var err error = nil

	for i := 0; i < 5; i++ {
		item, getErr := c.bfc.Get(key)

		if getErr == nil {
			return item, nil
		}

		err = getErr

		time.Sleep(15 * time.Millisecond)
	}

	fmt.Println(err)

	return nil, err
}

// Set sets a key. Will retry 3 times.
func (c Client) Set(item *memcache.Item) error {
	var err error = nil

	for i := 0; i < 5; i++ {
		err = c.bfc.Set(item)

		if err == nil {
			return err
		}

		time.Sleep(15 * time.Millisecond)
	}

	if err != nil {
		fmt.Println(err)
	}

	return err
}

// Cas does a compare and swap
func (c Client) Cas(item *memcache.Item) error {
	return c.bfc.CompareAndSwap(item)
}

// Delete removes a key.
func (c Client) Delete(key string) error {
	return c.bfc.Delete(key)
}

// Flush flushes the cache.
func (c Client) Flush() error {
	return c.bfc.FlushAll()
}

// ListAppend updates a key whose value is a list.
func (c Client) ListAppend(key, toAdd string) error {
	var err error = nil

	for i := 0; i < 1000; i++ {
		item, getErr := c.bfc.Get(key)

		if getErr != nil {
			err = getErr

			fmt.Println(fmt.Sprintf("Get miss: %s (lookup key), %s (addable).", key, toAdd))

			continue
		}

		var oldItems []string

		json.Unmarshal(item.Value, &oldItems)

		newItems := append(oldItems, toAdd)
		newItemsMarshalled, _ := json.Marshal(newItems)
		item.Value = newItemsMarshalled
		casErr := c.bfc.CompareAndSwap(item)

		if casErr == nil {
			err = nil

			break
		} else {
			err = casErr
		}
	}

	if err == nil {
		fmt.Println(fmt.Sprintf("Appended %s to %s", toAdd, key))
	} else {
		fmt.Println(fmt.Sprintf("Failed to CAS %s/%s", toAdd, key))
	}

	return err
}

// ListRemove updates a key whose value is a list.
func (c Client) ListRemove(key, toRemove string) error {
	var err error = nil

	for i := 0; i < 1000; i++ {
		item, getErr := c.bfc.Get(key)

		if getErr != nil {
			err = getErr

			fmt.Println(fmt.Sprintf("Get miss: %s (lookup key), %s (removable).", key, toRemove))

			continue
		}

		var items []string

		json.Unmarshal(item.Value, &items)

		for i, curr := range items {
			if curr == toRemove {
				items = append(items[:i], items[i+1:]...)
				break
			}
		}

		itemsMarshalled, _ := json.Marshal(items)
		item.Value = itemsMarshalled
		casErr := c.bfc.CompareAndSwap(item)

		if casErr == nil {
			err = nil

			break
		} else {
			err = casErr
		}
	}

	if err == nil {
		fmt.Println(fmt.Sprintf("Removed %s from %s", toRemove, key))
	} else {
		fmt.Println(fmt.Sprintf("Failed to CAS %s/%s", toRemove, key))
	}

	return err
}
