package lru

import "testing"

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
