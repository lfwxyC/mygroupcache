package geecache

import "container/list"

type LFUCache struct {
	maxBytes int64        // 允许使用的最大内存
	nBytes   int64        // 当前已使用的内存
	ll       []*list.List // 频率链表集合，用数组下标表示频率
	cache    map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为nil
	OnEvicted func(key string, value ByteView)
}

// list.Element.Value的数据类型
type lfuEntry struct {
	key   string // 淘汰队首节点时，用key从字典中删除对应的映射
	value ByteView
	freq  int // 频率
}

func NewLFUCache(maxBytes int64, onEvicted func(string, ByteView)) *LFUCache {
	return &LFUCache{
		maxBytes:  maxBytes,
		ll:        []*list.List{list.New()}, // 初始时只有频率为0的链表
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *LFUCache) Add(key string, value ByteView) {
	freq := 0
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		kv.value = value
		c.nBytes += int64(value.Len()) - int64(kv.value.Len())
		// 从原链表删除
		c.ll[kv.freq].Remove(ele)
		freq = kv.freq + 1
	} else {
		c.nBytes += int64(len(key)) + int64(value.Len())
	}
	// 加入新频率链表
	if len(c.ll) <= freq {
		c.ll = append(c.ll, list.New())
	}
	ele := c.ll[freq].PushFront(&lfuEntry{key: key, value: value, freq: freq})
	c.cache[key] = ele
	for c.maxBytes != 0 && c.nBytes > c.maxBytes {
		c.removeOldest()
	}
}

func (c *LFUCache) Get(key string) (value ByteView, ok bool) {
	if ele, ok := c.cache[key]; ok {
		kv := ele.Value.(*lfuEntry)
		// 从原链表删除
		c.ll[kv.freq].Remove(ele)
		// 加入新频率链表
		if len(c.ll) <= kv.freq+1 {
			c.ll = append(c.ll, list.New())
		}
		ele = c.ll[kv.freq].PushFront(&lfuEntry{key: key, value: value, freq: kv.freq + 1})
		return kv.value, true
	}
	return
}

func (c *LFUCache) removeOldest() {
	// 删除频率为0的链表的尾部节点
	ele := c.ll[0].Back()
	if ele != nil {
		kv := ele.Value.(*lruEntry)
		delete(c.cache, kv.key)
		c.nBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		c.ll[0].Remove(ele)
	}
}
