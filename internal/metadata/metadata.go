package metadata

import (
	"fmt"
	"image"
	"io"

	_ "github.com/gen2brain/avif"
)

// ImageMetadata contains basic image information
type ImageMetadata struct {
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	Format      string  `json:"format"`
	MimeType    string  `json:"mimeType"`
	AspectRatio float64 `json:"aspectRatio"`
	Dominant    struct {
		RGB  string `json:"rgb"`
		HSL  string `json:"hsl"`
		Name string `json:"name"`
	} `json:"dominant"`
	Colors []string `json:"colors"` // Palette extraction
	EXIF   map[string]interface{} `json:"exif,omitempty"`
}

// Get retrieves metadata from an image without loading the entire image into memory
func Get(r io.Reader) (ImageMetadata, error) {
	// Decode only the image config (header) which is much faster than decoding the whole image
	config, format, err := image.DecodeConfig(r)
	if err != nil {
		return ImageMetadata{}, fmt.Errorf("failed to decode image config: %w", err)
	}

	mimeType := fmt.Sprintf("image/%s", format)
	switch format {
	case "jpeg":
		mimeType = "image/jpeg"
	case "avif":
		mimeType = "image/avif"
	}

	return ImageMetadata{
		Width:    config.Width,
		Height:   config.Height,
		Format:   format,
		MimeType: mimeType,
	}, nil
}