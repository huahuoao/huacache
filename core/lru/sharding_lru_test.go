package lru

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func TestShardingLRU_Hash(t *testing.T) {
	sh, _ := NewShardingLRU(8, 8*1024)
	counts := make([]int, 8)
	for i := 0; i < 8000; i++ {
		n := sh.hash(randomString(64))
		counts[n]++
	}

	t.Log(counts)
}

func single(lru *ShardingLRU, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		lru.Set(randomString(100), randomString(10))
	}
}
func TestSetLRUSharding(t *testing.T) {
	lru, _ := NewShardingLRU(8, 8*1024*1024*1024)
	start := time.Now()
	wg := &sync.WaitGroup{}
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go single(lru, wg)
	}
	wg.Wait()
	fmt.Printf("%v", time.Since(start))
}
