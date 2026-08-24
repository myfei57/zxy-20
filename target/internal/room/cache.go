package room

import "sync"

// Cache is the in-memory view of room state served to the console pages.
type Cache struct {
	mu    sync.RWMutex
	rooms map[string]Room
}

// NewCache returns an empty room cache.
func NewCache() *Cache {
	return &Cache{rooms: make(map[string]Room)}
}

// Set stores the latest room state in memory.
func (c *Cache) Set(r Room) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.rooms[r.ID] = r
}

// Get returns the in-memory room state.
func (c *Cache) Get(id string) (Room, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	r, ok := c.rooms[id]
	return r, ok
}
