package lru

import (
	"container/list"
	"fmt"
)

// Cache is a LRU cache. It is not safe for concurrent access.
type Cache struct {
	maxBytes int64 //cache max byte limit
	nbytes   int64 //used bytes
	ll       *list.List
	cache    map[string]*list.Element
	// optional and executed when an entry is purged.
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

// Value use Len to count how many bytes it takes
type Value interface {
	Len() int
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

func (c *Cache) GetMemoryUsedSituation() (int64, int64, int) {
	return c.maxBytes, c.nbytes, c.ll.Len()
}
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) DeleteKey(key string) error {
	if ele, ok := c.cache[key]; ok {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		return nil
	}
	return fmt.Errorf("key unexited")
}

// Add adds a value to the cache.
func (c *Cache) Add(key string, value Value) error {
	if c.maxBytes != 0 && int64(len(key))+int64(value.Len()) > c.maxBytes {
		return fmt.Errorf("new item exceeds cache maximum limit")
	}

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}

	return nil
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.ll.Len()
}
func (c *Cache) Keys() ([]string, error) {
    keys := make([]string, 0, c.ll.Len())
    
    for e := c.ll.Front(); e != nil; e = e.Next() {
        kv, ok := e.Value.(*entry)
        if !ok {
            return nil, fmt.Errorf("invalid value type")
        }
        keys = append(keys, kv.key)
    }
    
    return keys, nil
}
