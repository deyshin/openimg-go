package cache

import (
	"testing"
)

func TestCache(t *testing.T) {
	cache := New()

	// Test setting and getting
	key := "test_key"
	data := []byte("test_data")

	cache.Set(key, data)
	got, found := cache.Get(key)

	if !found {
		t.Error("Expected to find key in cache")
	}

	if string(got) != string(data) {
		t.Errorf("Got %s, want %s", got, data)
	}

	// Test key not found
	_, found = cache.Get("nonexistent")
	if found {
		t.Error("Expected key to not be found")
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		width    int
		height   int
		quality  int
		format   string
		fit      string
		wantKey  string
	}{
		{
			name:    "basic key",
			url:     "http://example.com/image.jpg",
			width:   100,
			height:  100,
			quality: 80,
			format:  "jpeg",
			fit:     "cover",
			wantKey: "http://example.com/image.jpg_w100_h100_q80_fmtjpeg_fitcover",
		},
		{
			name:    "zero values",
			url:     "http://example.com/image.jpg",
			width:   0,
			height:  0,
			quality: 0,
			format:  "",
			fit:     "",
			wantKey: "http://example.com/image.jpg_w0_h0_q0_fmt_fit",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateKey(tt.url, tt.width, tt.height, tt.quality, tt.format, tt.fit)
			if got != tt.wantKey {
				t.Errorf("GenerateKey() = %v, want %v", got, tt.wantKey)
			}
		})
	}
}