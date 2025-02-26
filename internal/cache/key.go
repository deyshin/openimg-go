package cache

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

// GenerateKey generates a cache key from the image URL and transformation options
func GenerateKey(url string, width, height, quality int, format, fit string) string {
	// Create a unique key based on URL and transformation parameters
	key := fmt.Sprintf("%s_w%d_h%d_q%d_fmt%s_fit%s",
		url, width, height, quality, format, fit)

	// Hash the key to ensure safe characters and fixed length
	h := sha256.New()
	h.Write([]byte(key))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}