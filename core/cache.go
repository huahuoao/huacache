package huacache

import (
	"fmt"
	"sync"

	"github.com/huahuoao/huacache/core/lru"
)

type cache struct {
	mu         sync.Mutex
	lru        *lru.Cache
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	err := c.lru.Add(key, value)
	return err
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (c *cache) delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return fmt.Errorf("cache is uninitialized")
	}
	err := c.lru.DeleteKey(key)
	return err
}
