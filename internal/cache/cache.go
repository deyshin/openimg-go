package cache

import (
	"errors"
	"time"

	"github.com/coocood/freecache"
)

var ErrNotFound = errors.New("item not found in cache")

type Cache interface {
	Get(key string) ([]byte, error)
	Set(key string, value []byte) error
}

// Package-level constructor functions
func NewMemoryCache(sizeMB int, ttl time.Duration) Cache {
	sizeBytes := sizeMB * 1024 * 1024
	return &MemoryCache{
		cache: freecache.NewCache(sizeBytes),
	}
}

func NewNoopCache() Cache {
	return &NoopCache{}
}

func NewDiskCache(path string) Cache {
	return &DiskCache{
		basePath: path,
	}
}

// Options represents cache configuration
type Options struct {
	Type string // "memory", "disk", "redis", "s3", "none"
	Size int    // Size in MB for memory cache
	TTL  time.Duration
	Path string // Path for disk cache or URL for remote caches
}