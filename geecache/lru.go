package geecache

import (
	"container/list"
)

type LRUCache struct {
	maxBytes int64      // 允许使用的最大内存
	nBytes   int64      // 当前已使用的内存
	ll       *list.List // 双向链表
	cache    map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为nil
	OnEvicted func(key string, value ByteView)
}

// list.Element.Value的数据类型
type lruEntry struct {
	key   string // 淘汰队首节点时，用key从字典中删除对应的映射
	value ByteView
}

func NewLRUCache(maxBytes int64, onEvicted func(string, ByteView)) *LRUCache {
	return &LRUCache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *LRUCache) Add(key string, value ByteView) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lruEntry)
		kv.value = value
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		c.ll.MoveToFront(ele)
	} else {
		ele := c.ll.PushFront(&lruEntry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *LRUCache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*lruEntry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		c.ll.Remove(ele)
	}
}

func (c *LRUCache) Get(key string) (value ByteView, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*lruEntry)
		return kv.value, true
	}
	return
}

func (c *LRUCache) Len() int {
	return c.ll.Len()
}
