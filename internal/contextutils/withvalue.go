package contextutils

import "context"

// Key is string type
type Key string

const (
	// Params key
	Params = Key("params")
)

// Set context with value
func Set(ctx context.Context, key Key, val map[string]interface{}) context.Context {
	h, ok := ctx.Value(key).(map[string]interface{})

	if !ok {
		h = val
	} else {
		for k, v := range val {
			h[k] = v
		}
	}

	return context.WithValue(ctx, key, h)
}

// Get context value
func Get(ctx context.Context, key Key) map[string]interface{} {
	m, ok := ctx.Value(key).(map[string]interface{})

	if !ok {
		return map[string]interface{}{}
	}
	return m
}
