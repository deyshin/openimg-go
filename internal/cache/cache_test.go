package cache

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	cache := NewMemoryCache(100, time.Hour)

	// Test setting and getting
	key := "test_key"
	data := []byte("test_data")

	cache.Set(key, data)
	got, err := cache.Get(key)

	if err != nil {
		t.Error("Expected to find key in cache")
	}

	if string(got) != string(data) {
		t.Errorf("Got %s, want %s", got, data)
	}

	// Test key not found
	_, err = cache.Get("nonexistent")
	if err != ErrNotFound {
		t.Error("Expected key to not be found")
	}
}

func TestNoopCache(t *testing.T) {
	cache := NewNoopCache()

	// Test setting and getting
	key := "test_key"
	data := []byte("test_data")

	cache.Set(key, data)
	got, err := cache.Get(key)

	if err != ErrNotFound {
		t.Error("Expected key to not be found")
	}

	if got != nil {
		t.Errorf("Expected nil, got %s", got)
	}
}

func TestDiskCache(t *testing.T) {
	dir, err := ioutil.TempDir("", "diskcache")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	cache := NewDiskCache(dir)

	// Test setting and getting
	key := "test_key"
	data := []byte("test_data")

	cache.Set(key, data)
	got, err := cache.Get(key)

	if err != nil {
		t.Error("Expected to find key in cache")
	}

	if string(got) != string(data) {
		t.Errorf("Got %s, want %s", got, data)
	}

	// Test key not found
	_, err = cache.Get("nonexistent")
	if err != ErrNotFound {
		t.Error("Expected key to not be found")
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		width   int
		height  int
		quality int
		format  string
		fit     string
	}{
		{
			name:    "basic key",
			url:     "http://example.com/image.jpg",
			width:   100,
			height:  100,
			quality: 80,
			format:  "jpeg",
			fit:     "cover",
		},
		{
			name:    "zero values",
			url:     "http://example.com/image.jpg",
			width:   0,
			height:  0,
			quality: 0,
			format:  "",
			fit:     "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateKey(tt.url, tt.width, tt.height, tt.quality, tt.format, tt.fit)
			if got == "" {
				t.Error("GenerateKey() returned empty string")
			}
			// Check if the generated key is a valid base64 URL-encoded string
			if _, err := base64.URLEncoding.DecodeString(got); err != nil {
				t.Errorf("GenerateKey() returned invalid base64 URL-encoded string: %v", got)
			}
		})
	}
}
