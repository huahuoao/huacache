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
	cache, err := NewGroup(generateRandomString(5), MB*100)
	if err != nil {
		t.Fatalf("Failed to create cache group: %v", err)
	}
	if cache == nil {
		t.Fatal("Cache group is nil")
	}

	const numGoroutines = 100
	const numEntries = 1000000
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numEntries/numGoroutines; j++ {
				key := generateRandomString(100)
				cache.AddOrUpdate(key, ByteView{B: []byte("value1")})
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)
	t.Logf("Time taken: %s", elapsed)
	t.Logf("1000 sets Time taken: %s", elapsed/numEntries*1000) // Ensure this logs when using 'go test -v'
	t.Log(cache.GetStatus().ToString())
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
