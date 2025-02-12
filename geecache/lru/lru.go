package lru

import "container/list"

type Cache struct {
	maxBytes int64      // 允许使用的最大内存
	nBytes   int64      // 当前已使用的内存
	ll       *list.List // 双向链表
	cache    map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为nil
	OnEvicted func(key string, value Value)
}

// list.Element.Value的数据类型
type entry struct {
	key   string // 淘汰队首节点时，用key从字典中删除对应的映射
	value Value
}

type Value interface {
	Len() int // 值所占用的内存大小
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*entry)
		kv.value = value
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		c.ll.MoveToFront(ele)
	} else {
		ele := c.ll.PushFront(&entry{key: key, value: value})
		c.cache[key] = ele
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		c.ll.Remove(ele)
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

func (c *Cache) Len() int {
	return c.ll.Len()
}
