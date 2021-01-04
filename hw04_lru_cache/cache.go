package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear() // Очистить кэш
}

type lruCache struct {
	// Place your code here:
	// - capacity
	// - queue
	// - items
	_ cacheItem
}

type cacheItem struct {
	// Place your code here
}

func NewCache(_ int) Cache {
	return &lruCache{}
}

func (c *lruCache) Set(_ Key, _ interface{}) bool {
	return false
}

func (c *lruCache) Get(_ Key) (interface{}, bool) {
	return nil, false
}

func (c *lruCache) Clear() {

}
