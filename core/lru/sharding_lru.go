package lru

import (
	"errors"
	"github.com/spaolacci/murmur3"
)

type ShardingLRU struct {
	ShardingMap map[int]*Cache
	SliceNum    int
}
type String1 struct {
	str string
}

func (s String1) Len() int {
	return len(s.str)
}
func (sh *ShardingLRU) hash(s string) int {
	hashValue := murmur3.Sum32([]byte(s)) // 使用MurmurHash3
	return int(hashValue) % sh.SliceNum   // 取模
}
func (sh *ShardingLRU) GetLru(key string) *Cache {
	cache, exists := sh.ShardingMap[sh.hash(key)]
	if !exists {
		return nil
	}
	return cache
}
func NewShardingLRU(sliceNum int, maxBytes int64) (*ShardingLRU, error) {
	shardingMap := make(map[int]*Cache, sliceNum)
	if maxBytes%int64(sliceNum) != 0 {
		return nil, errors.New("maxBytes must be multiple of sliceNum")
	}
	for i := 0; i < sliceNum; i++ {
		shardingMap[i] = New(maxBytes/int64(sliceNum), nil)
	}
	return &ShardingLRU{
		ShardingMap: shardingMap, // 根据sliceNum初始化map大小
		SliceNum:    sliceNum,
	}, nil
}

func (sh *ShardingLRU) Set(key string, value string) {
	cache := sh.GetLru(key)
	if cache == nil {
		return
	}
	cache.Add(key, String1{str: value})
}
