package geecache

import (
	"sync"
)

type CacheType int

const (
	LRU = iota
	LFU
)

// Cache 淘汰策略
type Cache interface {
	Add(key string, value ByteView)
	Get(key string) (value ByteView, ok bool)
}

type cache struct {
	mu         sync.Mutex
	used       Cache
	cacheType  CacheType
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.used == nil {
		switch c.cacheType {
		case LRU:
			c.used = NewLRUCache(c.cacheBytes, nil)
		case LFU:
			c.used = NewLFUCache(c.cacheBytes, nil)
		}
	}
	c.used.Add(key, value)
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.used != nil {
		if v, ok := c.used.Get(key); ok {
			return v, true
		}
	}
	return
}
