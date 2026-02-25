package core

import "sync"

// Context carries state throughout a flow execution.
// It is shared across all nodes in a single flow run (FlowMemory),
// and is safe for concurrent access across goroutines.
type Context struct {
	mu     sync.RWMutex
	values map[string]any

	// FlowID is the identifier of the flow this context belongs to.
	FlowID string

	// TraceID is a unique identifier for this specific execution run.
	TraceID string
}

// NewContext creates a new flow execution context.
func NewContext(flowID, traceID string) *Context {
	return &Context{
		FlowID:  flowID,
		TraceID: traceID,
		values:  make(map[string]any),
	}
}

// Set stores a value in the context under the given key.
func (c *Context) Set(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = value
}

// Get retrieves a value from the context by key.
// Returns nil if the key does not exist.
func (c *Context) Get(key string) any {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.values[key]
}

// Has reports whether a key exists in the context.
func (c *Context) Has(key string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, ok := c.values[key]
	return ok
}

// Delete removes a value from the context.
func (c *Context) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.values, key)
}
