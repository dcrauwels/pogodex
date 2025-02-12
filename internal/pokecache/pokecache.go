package pokecache

import (
	"sync"
	"time"
)

// define entry struct
type cacheEntry struct {
	createdAt time.Time
	Val       []byte
}

// define main Cache struct
type Cache struct {
	Entry    map[string]cacheEntry
	mu       sync.RWMutex
	interval time.Duration
	stopChan chan struct{}
}

// function to create and return cache
func NewCache(interval int) *Cache {
	// sanity (want to check if interval < 0, not sure what to do yet)

	// turn interval int to timeInterval time.Duration
	timeInterval := time.Duration(interval) * time.Second

	// construct returned cache
	c := &Cache{
		Entry:    make(map[string]cacheEntry),
		interval: timeInterval,
		stopChan: make(chan struct{}),
	}

	// start separate goroutine for reapLoop()
	go c.reapLoop()

	// and return cache
	return c
}

// method to add cacheEntry to Cache
func (c *Cache) Add(key string, val []byte) {
	// mutex w lock
	c.mu.Lock()
	defer c.mu.Unlock()

	// add entry
	c.Entry[key] = cacheEntry{
		createdAt: time.Now(),
		Val:       val,
	}
}

// method to retrieve cacheEntry from Cache, return false if key not found
func (c *Cache) Get(key string) ([]byte, bool) {
	// mutex r lock
	c.mu.RLock()
	defer c.mu.RUnlock()

	// check if key in Cache.Entry slice
	value, ok := c.Entry[key]

	// init zero value return
	var val []byte

	// update return value
	if ok {
		val = value.Val
	}

	// return
	return val, ok
}

// method to periodically call reap()
func (c *Cache) reapLoop() {
	// create ticker for periodic reap() calls
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	// endless loop
	for {
		select {
		case <-ticker.C:
			c.reap()
		case <-c.stopChan: //check for stop signal
			return
		}
	}
}

// method to delete Entry from Cache, called periodically by reapLoop()
func (c *Cache) reap() {

	c.mu.Lock()
	defer c.mu.Unlock()
	now := time.Now()
	for key, value := range c.Entry {
		if now.Sub(value.createdAt) > c.interval {
			delete(c.Entry, key)
			//fmt.Printf("Removed key '%s' from cache.\n", key)
		}
	}

}

// method to gracefully stop reap loop
func (c *Cache) Stop() {
	close(c.stopChan)
}
