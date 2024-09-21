package huacache

import (
	"crypto/rand"
	"math/big"
	"sync"
	"testing"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func generateRandomString(length int) string {
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return ""
		}
		result[i] = letterBytes[num.Int64()]
	}
	return string(result)
}
func TestAddOrUpdate(t *testing.T) {
	cache, _ := NewGroup("test", MB*100)
	cache.AddOrUpdate("key1", ByteView{B: []byte("value1")})
	v, _ := cache.Get("key1")
	if string(v.B) != "value1" {
		t.Fatalf("cache get failed")
	}
	v, _ = cache.Get("key not exist")
	if string(v.B) != "" {
		t.Fatalf("cache get empty key failed")
	}
	cache.AddOrUpdate("key1", ByteView{B: []byte("value2")})
	v, _ = cache.Get("key1")
	if string(v.B) != "value2" {
		t.Fatalf("cache update failed")
	}
}

func TestConcurrentSetKey(t *testing.T) {
	cache, err := NewGroup(generateRandomString(5), MB*10000)
	if err != nil {
		t.Fatalf("Failed to create cache group: %v", err)
	}
	if cache == nil {
		t.Fatal("Cache group is nil")
	}
	const numGoroutines = 8
	var wg sync.WaitGroup
	start := time.Now()
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				key := generateRandomString(10)
				err := cache.AddOrUpdate(key, ByteView{B: []byte(key)})
				if err != nil {
					return
				}
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("Time taken: %s", elapsed)
}

func TestDelete(t *testing.T) {
	cache, _ := NewGroup(generateRandomString(5), MB*100)
	cache.AddOrUpdate("key1", ByteView{B: []byte("value1")})
	v, _ := cache.Get("key1")
	if string(v.B) != "value1" {
		t.Fatalf("cache get failed")
	}
	cache.Delete("key1")
	v, _ = cache.Get("key1")
	if string(v.B) != "" {
		t.Fatalf("cache remove failed")
	}
}
