package cache

import (
	"context"
	"encoding/json"
	"sync"
)

// Map cache
type Map struct {
	sync.RWMutex
	c map[string]interface{}
}

//NewMap is contructor.
func NewMap() *Map {
	return &Map{
		c: make(map[string]interface{}),
	}
}

// Get key
func (c *Map) Get(ctx context.Context, key string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()

	v, _ := c.c[key]
	return v, nil
}

// Delete key
func (c *Map) Delete(ctx context.Context, key string) (interface{}, error) {
	v, _ := c.Get(ctx, key)

	c.Lock()
	defer c.Unlock()

	delete(c.c, key)
	return v, nil
}

// Set key val
func (c *Map) Set(ctx context.Context, key string, val interface{}, ttl ...int) error {
	c.Lock()
	defer c.Unlock()

	c.c[key] = val
	return nil
}

// Inc a key.
func (c *Map) Inc(ctx context.Context, key string) (int, error) {
	c.Lock()
	defer c.Unlock()

	v, _ := c.c[key].(int)
	v++
	c.c[key] = v

	return v, nil
}

func (c *Map) String() string {
	c.RLock()
	defer c.RUnlock()

	d, _ := json.Marshal(c.c)
	return string(d)
}
