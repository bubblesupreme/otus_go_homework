package hw04_lru_cache //nolint:golint,stylecheck

import log "github.com/sirupsen/logrus"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*listItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	if capacity < 1 {
		log.Warn("capacity is less than 1, cache will always be empty")
	}
	return &lruCache{capacity, NewList(), make(map[Key]*listItem, capacity)}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	if c.capacity < 1 {
		return false
	}

	newItem := cacheItem{key, value}
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		c.queue.Front().Value = newItem
		return true
	}
	if c.capacity <= c.queue.Len() {
		delete(c.items, c.queue.Back().Value.(cacheItem).key)
		c.queue.Remove(c.queue.Back())
	}
	c.items[key] = c.queue.PushFront(newItem)
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	if v, ok := c.items[key]; ok {
		c.queue.MoveToFront(v)
		return v.Value.(cacheItem).value, true
	}
	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*listItem, c.capacity)
}
