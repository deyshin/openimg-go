package validate

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	MaxWidth  = 2000
	MaxHeight = 2000
	MinWidth  = 1
	MinHeight = 1
)

// ImageOptions validates transformation parameters
func ImageOptions(width, height, quality int, format, fit string) error {
	if width != 0 && (width < MinWidth || width > MaxWidth) {
		return fmt.Errorf("width must be between %d and %d", MinWidth, MaxWidth)
	}
	if height != 0 && (height < MinHeight || height > MaxHeight) {
		return fmt.Errorf("height must be between %d and %d", MinHeight, MaxHeight)
	}
	if quality != 0 && (quality < 1 || quality > 100) {
		return fmt.Errorf("quality must be between 1 and 100")
	}
	if format != "" && format != "jpeg" && format != "png" {
		return fmt.Errorf("format must be jpeg or png")
	}
	if fit != "" && fit != "cover" && fit != "contain" {
		return fmt.Errorf("fit must be cover or contain")
	}
	return nil
}

// URL validates the source image URL
func URL(rawURL string) error {
	if rawURL == "" {
		return fmt.Errorf("URL is required")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return fmt.Errorf("URL scheme must be http or https")
	}

	return nil
}