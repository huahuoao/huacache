package huacache

import (
	"fmt"

	"github.com/huahuoao/huacache/core/lru"
)

type cache struct {
	lru        *lru.ShardingLRU
	cacheBytes int64
}

func (c *cache) add(key string, value ByteView) error {
	err := c.lru.GetLru(key).Add(key, value)
	return err
}

func (c *cache) get(key string) (value ByteView, ok bool) {
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.GetLru(key).Get(key); ok {
		return v.(ByteView), ok
	}

	return
}

func (c *cache) delete(key string) error {
	if c.lru == nil {
		return fmt.Errorf("cache is uninitialized")
	}
	err := c.lru.GetLru(key).DeleteKey(key)
	return err
}
