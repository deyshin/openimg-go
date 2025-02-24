package cache

import (
	"fmt"
	"sync"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type ImageCache struct {
	cache *gocache.Cache
	mu    sync.RWMutex
}

// New creates a new ImageCache with default expiration of 1 hour and cleanup every 2 hours
func New() *ImageCache {
	return &ImageCache{
		cache: gocache.New(1*time.Hour, 2*time.Hour),
	}
}

// Get retrieves an image from the cache
func (c *ImageCache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if data, found := c.cache.Get(key); found {
		return data.([]byte), true
	}
	return nil, false
}

// Set stores an image in the cache
func (c *ImageCache) Set(key string, data []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.cache.Set(key, data, gocache.DefaultExpiration)
}

// GenerateKey generates a cache key from the image URL and transformation options
func GenerateKey(url string, width, height, quality int, format, fit string) string {
	return fmt.Sprintf("%s_w%d_h%d_q%d_fmt%s_fit%s",
		url, width, height, quality, format, fit)
}