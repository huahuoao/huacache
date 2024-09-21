package lru

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGet(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key", String("value"))
	v, ok := lru.Get("key")
	if !ok {
		t.Fatalf("cache get failed")
	}
	if ok && v != String("value") {
		t.Fatalf("cache get failed")
	}
}

func TestMemoryElimination(t *testing.T) {
	cap := len("key1" + "value1")
	lru := New(int64(cap), nil)
	lru.Add("key1", String("value1"))
	lru.Add("key2", String("value2")) //now key1 is eliminating
	_, ok := lru.Get("key1")
	if ok {
		t.Fatalf("memory eliminate failed")
	}
	_, ok = lru.Get("key2")
	if !ok {
		t.Fatalf("memory eliminate failed")
	}
}

func TestRemoveKey(t *testing.T) {
	lru := New(int64(0), nil)
	lru.Add("key1", String("value1"))
	lru.DeleteKey("key1")
	_, ok := lru.Get("key1")
	if ok {
		t.Fatalf("memory remove key failed")
	}
}

func TestOnEvicted(t *testing.T) {
	now_string := "key"
	callback := func(key string, value Value) {
		now_string += key
	}
	lru := New(int64(0), callback)
	lru.Add("new_key", String("value"))
	expect := "keynew_key"
	lru.DeleteKey("new_key")
	if now_string != expect {
		t.Fatalf("Call OnEvivted failed")
	}
}
func single1(lru *Cache, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < 1000; i++ {
		lru.Add(randomString(100), String(randomString(10)))
	}
}
func TestSet(t *testing.T) {
	lru := New(8*1024*1024*1024, nil)
	start := time.Now()
	var wg sync.WaitGroup
	wg.Add(8)
	for i := 0; i < 8; i++ {
		go single1(lru, &wg)
	}
	wg.Wait()
	fmt.Printf("%v", time.Since(start))
}
