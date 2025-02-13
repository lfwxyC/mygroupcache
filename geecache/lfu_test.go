package geecache

import "testing"

func TestGetLFU(t *testing.T) {
	lfu := NewLFUCache(int64(0), nil)
	lfu.Add("key1", ByteView{b: cloneBytes([]byte("0418"))})
	if v, ok := lfu.Get("key1"); !ok || v.String() != "0418" {
		t.Fatal("cache hit key1=0224 failed")
	}
	if _, ok := lfu.Get("key2"); ok {
		t.Fatal("cache miss key2 failed")
	}
}

func TestLFURemoveOldest(t *testing.T) {
	k1, k2, k3 := "key1", "key2", "k3"
	v1, v2, v3 := "value1", "value2", "v3"
	cap := len(k1 + k2 + v1 + v2)
	lru := NewLRUCache(int64(cap), nil)
	lru.Add(k1, ByteView{b: cloneBytes([]byte(v1))})
	lru.Add(k2, ByteView{b: cloneBytes([]byte(v2))})
	lru.Get(k1)
	lru.Add(k3, ByteView{b: cloneBytes([]byte(v3))})

	if _, ok := lru.Get(k2); ok {
		t.Fatal("RemoveOldest key2 failed")
	}
}
