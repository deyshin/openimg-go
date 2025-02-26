package cache

type NoopCache struct{}

func (c *NoopCache) Get(key string) ([]byte, error) {
    return nil, ErrNotFound
}

func (c *NoopCache) Set(key string, value []byte) error {
    return nil
}