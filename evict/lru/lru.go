package lru

import (
	"container/list"
	"fmt"
)

type Cache struct {
	// 最大容量
	maxBytes int64
	// 已存容量
	nBytes int64
	// 双向链表用于存储最近使用顺序
	ll *list.List
	// 快速索引缓存
	hashMap map[string]*list.Element
	// 淘汰回调方法
	// todo 这里是否可以设置为列表
	OnEvicted func(key string, value Value)
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int64
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		hashMap:   make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

func (c *Cache) Get(key string) (value Value, err error) {
	if ele, ok := c.hashMap[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, nil
	}
	return nil, fmt.Errorf("not found")
}

func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.hashMap[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nBytes += value.Len() - kv.value.Len()
		kv.value = value
	} else {
		ele := c.ll.PushFront(&entry{key, value})
		c.hashMap[key] = ele
		c.nBytes += int64(len(key)) + value.Len()
	}
	for c.maxBytes != 0 && c.maxBytes < c.nBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		kv := ele.Value.(*entry)
		c.nBytes -= int64(len(kv.key)) + kv.value.Len()
		c.ll.Remove(ele)
		delete(c.hashMap, kv.key)
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Len() int64 {
	return int64(c.ll.Len())
}
