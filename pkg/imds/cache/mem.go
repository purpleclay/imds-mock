/*
Copyright (c) 2022 Purple Clay

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package cache

import "sync"

// MemCache defines a lightweight in-memory cache that is thread safe
type MemCache struct {
	items map[string]string
	mu    sync.RWMutex
}

// New will generate an return a new empty in-memory cache
func New() *MemCache {
	return &MemCache{
		items: map[string]string{},
	}
}

// Set a value within the cache using the given cache key. If the value already
// exists, it will be overwritten
func (c *MemCache) Set(key string, value string) {
	c.mu.Lock()
	c.items[key] = value
	c.mu.Unlock()
}

// Get returns an item from the cache using the given cache key. A flag is also returned
// indicating whether the item exists. If no item exists, an empty string will be returned
func (c *MemCache) Get(key string) (string, bool) {
	c.mu.RLock()
	item, exists := c.items[key]
	c.mu.RUnlock()

	return item, exists
}

// Remove an item from the cache using the given cache key. Nothing will happen if no
// item exists
func (c *MemCache) Remove(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}
