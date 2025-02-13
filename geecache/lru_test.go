package geecache

import (
	"reflect"
	"testing"
)

type String string

func (s String) Len() int {
	return len(s)
}

func TestGetLRU(t *testing.T) {
	lru := NewLRUCache(int64(0), nil)
	lru.Add("key1", ByteView{b: cloneBytes([]byte("0418"))})
	if v, ok := lru.Get("key1"); !ok || v.String() != "0418" {
		t.Fatal("cache hit key1=0224 failed")
	}
	if _, ok := lru.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestLRURemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewLRUCache(int64(cap), nil)
	lru.Add(k1, ByteView{b: cloneBytes([]byte(v1))})
	lru.Add(k2, ByteView{b: cloneBytes([]byte(v2))})
	lru.Add(k3, ByteView{b: cloneBytes([]byte(v3))})

	if _, ok := lru.Get(k1); ok {
		t.Fatal("RemoveOldest key1 failed")
	}
}

func TestOnEvicted(t *testing.T) {
	keys := make([]string, 0)
	callback := func(key string, value ByteView) {
		keys = append(keys, key)
	}
	lru := NewLRUCache(int64(10), callback)
	lru.Add("key1", ByteView{b: cloneBytes([]byte("0418"))})
	lru.Add("k2", ByteView{b: cloneBytes([]byte("v1"))})
	lru.Add("k3", ByteView{b: cloneBytes([]byte("v2"))})
	lru.Add("k4", ByteView{b: cloneBytes([]byte("v3"))})

	expect := []string{"key1", "k2"}
	if !reflect.DeepEqual(keys, expect) {
		t.Fatalf("Call OnEvicted failed, keys = %s", keys)
	}
}
