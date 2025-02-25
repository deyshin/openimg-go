package transform

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
	"github.com/gen2brain/avif"
	"github.com/gen2brain/webp"
)

// Options represents image transformation options
type Options struct {
	Width   int
	Height  int
	Format  string // "jpeg", "jpg", "png", "avif", or "webp"
	Quality int    // 1-100, for JPEG, AVIF and WebP
	Fit     string // "cover" or "contain"
	Blur    float64 // Gaussian blur sigma
	Sharpen float64 // Sharpening intensity
	Rotate  float64 // Rotation angle in degrees
	Flip    string  // "horizontal", "vertical"
	Grayscale bool
	Blurhash bool // Generate blurhash string
	Smart    bool // Enable content-aware cropping
}

// PlaceholderOptions represents options for generating image placeholders
type PlaceholderOptions struct {
	Width   int
	Height  int
	Quality int // 1-100, lower means smaller size
}

const (
	FormatJPEG = "jpeg"
	FormatJPG  = "jpg"
	FormatPNG  = "png"
	FormatAVIF = "avif"
	FormatWEBP = "webp"
)

// Add AVIF-specific constants
const (
	DefaultAVIFQuality = 85
	DefaultAVIFSpeed = 8  // Balance between speed and compression
)

// Transform applies the specified transformations to an image
func Transform(img image.Image, opts Options) ([]byte, error) {
	// Apply resizing if needed
	if opts.Width > 0 || opts.Height > 0 {
		switch opts.Fit {
		case "cover":
			img = imaging.Fill(img, opts.Width, opts.Height, imaging.Center, imaging.Lanczos)
		case "contain":
			img = imaging.Fit(img, opts.Width, opts.Height, imaging.Lanczos)
		default:
			img = imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
		}
	}

	// Encode the image
	buf := new(bytes.Buffer)
	switch opts.Format {
	case "png":
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
	case "avif":
		quality := opts.Quality
		if quality == 0 {
			quality = DefaultAVIFQuality
		}
		// AVIF quality must be between 0 and 63
		quality = quality * 63 / 100  // Convert from 0-100 scale to 0-63 scale
		if err := avif.Encode(buf, img, avif.Options{
			Quality: quality,
			Speed:   DefaultAVIFSpeed,
		}); err != nil {
			return nil, fmt.Errorf("failed to encode AVIF: %w", err)
		}
	case "webp":
		quality := opts.Quality
		if quality == 0 {
			quality = 85
		}
		if err := webp.Encode(buf, img, webp.Options{
			Lossless: false,
			Quality:  quality,
		}); err != nil {
			return nil, err
		}
	case "jpg", "jpeg":
		quality := opts.Quality
		if quality == 0 {
			quality = 85
		}
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}
	default:
		// Default to PNG
		if err := png.Encode(buf, img); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

// GeneratePlaceholder creates a low-quality base64 placeholder
func GeneratePlaceholder(img image.Image, opts PlaceholderOptions) (string, error) {
	// Default values
	if opts.Width == 0 {
		opts.Width = 40 // very small width for placeholder
	}
	if opts.Height == 0 {
		bounds := img.Bounds()
		ratio := float64(bounds.Dy()) / float64(bounds.Dx())
		opts.Height = int(float64(opts.Width) * ratio)
	}
	if opts.Quality == 0 {
		opts.Quality = 20 // low quality for small size
	}

	// Resize image to small dimensions
	resized := imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)

	// Encode to JPEG with low quality
	buf := new(bytes.Buffer)
	if err := jpeg.Encode(buf, resized, &jpeg.Options{Quality: opts.Quality}); err != nil {
		return "", err
	}

	// Convert to base64
	b64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	return "data:image/jpeg;base64," + b64, nil
}

type ColorOptions struct {
	Brightness float64
	Contrast   float64
	Saturation float64
	Hue        float64
}

func AdjustColors(img image.Image, opts ColorOptions) image.Image {
	// Implement color adjustments
	return img
}

type ProgressiveOptions struct {
	Quality []int  // Multiple quality steps
	Sizes   []int  // Multiple size steps
	Format  string // Output format
}

func GenerateProgressiveImages(img image.Image, opts ProgressiveOptions) ([][]byte, error) {
	// Generate multiple versions for progressive loading
	return nil, nil
}