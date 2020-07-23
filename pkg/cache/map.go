package cache

import (
	"context"
	"encoding/json"
)

// Map cache
type Map map[string]interface{}

// Get key
func (c Map) Get(ctx context.Context, key string) (interface{}, error) {
	v, _ := c[key]
	return v, nil
}

// Delete key
func (c Map) Delete(ctx context.Context, key string) (interface{}, error) {
	v, _ := c.Get(ctx, key)
	delete(c, key)
	return v, nil
}

// Set key val
func (c Map) Set(ctx context.Context, key string, val interface{}, ttl ...int) error {
	c[key] = val
	return nil
}

// Inc a key.
func (c Map) Inc(ctx context.Context, key string) (int, error) {
	v, _ := c[key].(int)
	v++
	c[key] = v

	return v, nil
}

func (c Map) String() string {
	d, _ := json.Marshal(c)
	return string(d)
}
