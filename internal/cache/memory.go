package cache

import (
	"github.com/coocood/freecache"
)

type MemoryCache struct {
	cache *freecache.Cache
}

func (c *MemoryCache) Get(key string) ([]byte, error) {
	data, err := c.cache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *MemoryCache) Set(key string, value []byte) error {
	return c.cache.Set([]byte(key), value, 3600)
}