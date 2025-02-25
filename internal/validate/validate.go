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

var ValidFitModes = []string{
	"cover",
	"contain",
	"fill",
	"inside",
	"outside",
}

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
	if format != "" && !isValidFormat(format) {
		return fmt.Errorf("format must be one of: jpeg, jpg, png, avif, webp")
	}
	if fit != "" && !contains(ValidFitModes, fit) {
		return fmt.Errorf("fit must be one of: %v", ValidFitModes)
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

func contains(slice []string, item string) bool {
	for _, i := range slice {
		if i == item {
			return true
		}
	}
	return false
}

func isValidFormat(format string) bool {
	validFormats := []string{"jpeg", "jpg", "png", "avif", "webp"}
	return contains(validFormats, format)
}