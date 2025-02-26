package cache

import (
	"io"
	"os"
	"path/filepath"
)

type DiskCache struct {
	basePath string
}

func (c *DiskCache) Get(key string) ([]byte, error) {
	filePath := filepath.Join(c.basePath, key)
	file, err := os.Open(filePath)
	if os.IsNotExist(err) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (c *DiskCache) Set(key string, value []byte) error {
	filePath := filepath.Join(c.basePath, key)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(value)
	if err != nil {
		return err
	}
	return nil
}
