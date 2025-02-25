package metadata

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"testing"

	_ "image/jpeg"
	_ "image/png"

	"github.com/gen2brain/avif"
)

// createTestImage creates a simple test image with the given dimensions
func createTestImage(width, height int, format string) []byte {
	if width == 0 || height == 0 {
		return []byte("invalid image data")
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	buf := new(bytes.Buffer)

	switch format {
	case "png":
		png.Encode(buf, img)
	case "invalid":
		return []byte("invalid image data")
	case "avif":
		err := avif.Encode(buf, img, avif.Options{
			Quality: 60,
			Speed:   8,
		})
		if err != nil {
			// If AVIF encoding fails, return invalid image data
			return []byte("invalid image data")
		}
	default:
		jpeg.Encode(buf, img, nil)
	}

	return buf.Bytes()
}

func TestGet(t *testing.T) {
	tests := []struct {
		name       string
		width      int
		height     int
		format     string
		wantWidth  int
		wantHeight int
		wantFormat string
		wantMime   string
		wantErr    bool
	}{
		{
			name:       "jpeg image",
			width:      800,
			height:     600,
			format:     "jpeg",
			wantWidth:  800,
			wantHeight: 600,
			wantFormat: "jpeg",
			wantMime:   "image/jpeg",
			wantErr:    false,
		},
		{
			name:       "png image",
			width:      400,
			height:     300,
			format:     "png",
			wantWidth:  400,
			wantHeight: 300,
			wantFormat: "png",
			wantMime:   "image/png",
			wantErr:    false,
		},
		{
			name:    "invalid image data",
			width:   100,
			height:  100,
			format:  "invalid",
			wantErr: true,
		},
		{
			name:    "zero dimensions",
			width:   0,
			height:  0,
			format:  "jpeg",
			wantErr: true,
		},
		{
			name:       "avif image",
			width:      640,
			height:     480,
			format:     "avif",
			wantWidth:  640,
			wantHeight: 480,
			wantFormat: "avif",
			wantMime:   "image/avif",
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			imgData := createTestImage(tt.width, tt.height, tt.format)
			got, err := Get(bytes.NewReader(imgData))

			if (err != nil) != tt.wantErr {
				t.Errorf("error: = %v, wantErr: %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Width != tt.wantWidth {
					t.Errorf("Get() width = %v, want %v", got.Width, tt.wantWidth)
				}
				if got.Height != tt.wantHeight {
					t.Errorf("Get() height = %v, want %v", got.Height, tt.wantHeight)
				}
				if got.Format != tt.wantFormat {
					t.Errorf("Get() format = %v, want %v", got.Format, tt.wantFormat)
				}
				if got.MimeType != tt.wantMime {
					t.Errorf("Get() mimeType = %v, want %v", got.MimeType, tt.wantMime)
				}
			}
		})
	}
}