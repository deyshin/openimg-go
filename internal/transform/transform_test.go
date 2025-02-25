package transform

import (
	"bytes"
	"image"
	"image/color"
	"strings"
	"testing"

	"encoding/base64"

	"github.com/stretchr/testify/require"
)

// createTestImage creates a simple test image with the given dimensions
func createTestImage(width, height int) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Add some content to the image
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x % 256),
				G: uint8(y % 256),
				B: 100,
				A: 255,
			})
		}
	}
	return img
}

func TestTransform(t *testing.T) {
	tests := []struct {
		name    string
		img     image.Image
		opts    Options
		wantW   int
		wantH   int
		wantFmt string
		wantErr bool
	}{
		{
			name:    "resize to specific dimensions",
			img:     createTestImage(800, 600),
			opts:    Options{Width: 400, Height: 300, Format: "jpeg"},
			wantW:   400,
			wantH:   300,
			wantFmt: "jpeg",
			wantErr: false,
		},
		{
			name:    "convert to PNG",
			img:     createTestImage(100, 100),
			opts:    Options{Format: "png"},
			wantW:   100,
			wantH:   100,
			wantFmt: "png",
			wantErr: false,
		},
		{
			name: "resize with cover fit",
			img:  createTestImage(800, 600),
			opts: Options{
				Width:  400,
				Height: 400,
				Fit:    "cover",
				Format: "jpeg",
			},
			wantW:   400,
			wantH:   400,
			wantFmt: "jpeg",
			wantErr: false,
		},
		{
			name: "resize with contain fit",
			img:  createTestImage(800, 600),
			opts: Options{
				Width:  400,
				Height: 400,
				Fit:    "contain",
				Format: "jpeg",
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "jpeg",
			wantErr: false,
		},
		{
			name: "convert to jpg",
			img:  createTestImage(400, 300),
			opts: Options{
				Format:  "jpg",
				Quality: 85,
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "jpeg",
		},
		{
			name: "convert to avif",
			img:  createTestImage(400, 300),
			opts: Options{
				Format:  "avif",
				Quality: 85,
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "avif",
		},
		{
			name: "convert to webp",
			img:  createTestImage(400, 300),
			opts: Options{
				Format:  "webp",
				Quality: 85,
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "webp",
		},
		{
			name: "convert to avif with default quality",
			img:  createTestImage(400, 300),
			opts: Options{
				Format: "avif",
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "avif",
			wantErr: false,
		},
		{
			name: "convert to avif with custom quality",
			img:  createTestImage(400, 300),
			opts: Options{
				Format:  "avif",
				Quality: 60,
			},
			wantW:   400,
			wantH:   300,
			wantFmt: "avif",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transform(tt.img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Transform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			// Decode the result to verify dimensions and format
			resultImg, format, err := image.Decode(bytes.NewReader(got))
			if err != nil {
				t.Errorf("Failed to decode transformed image: %v", err)
				return
			}

			bounds := resultImg.Bounds()
			gotW := bounds.Dx()
			gotH := bounds.Dy()

			if gotW != tt.wantW || gotH != tt.wantH {
				t.Errorf("Transform() dimensions = %dx%d, want %dx%d",
					gotW, gotH, tt.wantW, tt.wantH)
			}

			if format != tt.wantFmt {
				t.Errorf("Transform() format = %v, want %v", format, tt.wantFmt)
			}
		})
	}
}

func TestTransform_Quality(t *testing.T) {
	img := createTestImage(400, 300)

	tests := []struct {
		name     string
		quality  int
		wantSize string // rough expectation: "smaller" or "larger"
	}{
		{
			name:     "high quality",
			quality:  100,
			wantSize: "larger",
		},
		{
			name:     "low quality",
			quality:  20,
			wantSize: "smaller",
		},
	}

	var lastSize int
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Transform(img, Options{
				Quality: tt.quality,
				Format:  "jpeg",
			})
			if err != nil {
				t.Errorf("Transform() error = %v", err)
				return
			}

			currentSize := len(got)
			if lastSize > 0 {
				if tt.wantSize == "smaller" && currentSize >= lastSize {
					t.Errorf("Expected size to be smaller than %d, got %d", lastSize, currentSize)
				}
			}
			lastSize = currentSize
		})
	}
}

func TestGeneratePlaceholder(t *testing.T) {
	// Create a test image
	img := createTestImage(800, 600)

	tests := []struct {
		name    string
		opts    PlaceholderOptions
		wantErr bool
	}{
		{
			name: "default options",
			opts: PlaceholderOptions{},
		},
		{
			name: "custom size",
			opts: PlaceholderOptions{
				Width:   20,
				Height:  15,
				Quality: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePlaceholder(img, tt.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePlaceholder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Basic validation of the data URL
			if !strings.HasPrefix(got, "data:image/jpeg;base64,") {
				t.Errorf("GeneratePlaceholder() result doesn't start with data URL prefix")
			}

			// Decode base64 to verify it's valid
			b64Data := strings.TrimPrefix(got, "data:image/jpeg;base64,")
			_, err = base64.StdEncoding.DecodeString(b64Data)
			if err != nil {
				t.Errorf("GeneratePlaceholder() produced invalid base64: %v", err)
			}
		})
	}
}

func TestAVIFQualityImpactsSize(t *testing.T) {
	img := createTestImage(400, 300)

	highQuality, err := Transform(img, Options{
		Format:  "avif",
		Quality: 90,
	})
	if err != nil {
		t.Fatalf("Failed to transform high quality: %v", err)
	}

	lowQuality, err := Transform(img, Options{
		Format:  "avif",
		Quality: 30,
	})
	if err != nil {
		t.Fatalf("Failed to transform low quality: %v", err)
	}

	if len(lowQuality) >= len(highQuality) {
		t.Errorf("Expected low quality AVIF to be smaller than high quality. Low: %d bytes, High: %d bytes",
			len(lowQuality), len(highQuality))
	}
}

func TestTransform_ConvertToAVIF(t *testing.T) {
	t.Run("convert_to_avif", func(t *testing.T) {
		result, err := Transform(createTestImage(100, 100), Options{
			Format: "avif",
			Width:  100,
			Height: 100,
		})
		require.NoError(t, err)

		// Verify the result is actually AVIF
		decoded, format, err := image.Decode(bytes.NewReader(result))
		require.NoError(t, err)
		require.NotNil(t, decoded)
		require.Equal(t, "avif", format)
	})
}