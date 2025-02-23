package transform

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"

	"github.com/disintegration/imaging"
)

// Options represents image transformation options
type Options struct {
	Width   int
	Height  int
	Format  string // "jpeg" or "png"
	Quality int    // 1-100, only for JPEG
	Fit     string // "cover" or "contain"
}

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
	default:
		// Default to JPEG
		quality := opts.Quality
		if quality == 0 {
			quality = 85
		}
		if err := jpeg.Encode(buf, img, &jpeg.Options{Quality: quality}); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}