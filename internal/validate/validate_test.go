package validate

import "testing"

func TestImageOptions(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		quality int
		format  string
		fit     string
		wantErr bool
	}{
		{"valid options", 800, 600, 80, "jpeg", "cover", false},
		{"zero dimensions", 0, 0, 80, "jpeg", "", false},
		{"width too large", MaxWidth + 1, 600, 80, "jpeg", "", true},
		{"height too small", 800, 0, 80, "jpeg", "", false},
		{"invalid quality", 800, 600, 101, "jpeg", "", true},
		{"invalid format", 800, 600, 80, "gif", "", true},
		{"invalid fit", 800, 600, 80, "jpeg", "stretch", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ImageOptions(tt.width, tt.height, tt.quality, tt.format, tt.fit)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImageOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestURL(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{"valid http", "http://example.com/image.jpg", false},
		{"valid https", "https://example.com/image.jpg", false},
		{"empty url", "", true},
		{"invalid scheme", "ftp://example.com/image.jpg", true},
		{"invalid format", "not-a-url", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := URL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("URL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}