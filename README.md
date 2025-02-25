# OpenImg Go

A Go image processing service that provides image transformation, metadata retrieval, and placeholder generation with caching support.

## Features

- **Image Transformation**
  - Resize with width/height parameters
  - Format conversion (JPEG, PNG, AVIF, WebP)
  - High-quality AVIF compression
  - Quality control for lossy formats
  - Multiple fit modes (cover, contain)
- **Image Metadata**
  - Dimensions (width, height)
  - Format detection
  - MIME type
- **Placeholder Generation**
  - Base64 encoded low-quality previews
  - Configurable dimensions and quality
- **Performance**
  - In-memory caching
  - Efficient metadata extraction
  - Optimized image processing

## Getting Started

### Prerequisites

Install Go (1.21 or later)

System dependencies:
- macOS:

  ```bash
  brew install aom  # Required for AVIF support
  ```

- Ubuntu/Debian:

  ```bash
  sudo apt-get install libaom-dev  # Required for AVIF support
  ```

### Installation

1. Clone the repository

   ```bash
   git clone https://github.com/deyshin/openimg-go
   cd openimg-go
   ```

2. Install dependencies

   ```bash
   go mod download
   ```

3. Install Air for hot reloading (optional)

   ```bash
   go install github.com/cosmtrek/air@latest
   ```

### Running the Server

Run with hot reloading:

  ```bash
  air
  ```

Or run normally:

  ```bash
  go run main.go
  ```

## API Endpoints

### Images
```
GET /api/image?url=<image_url>&w=<width>&h=<height>&fmt=<format>&q=<quality>&fit=<fit>
http://localhost:8080/api/image?url=https://example.com/image.jpg&w=800&h=600&fmt=jpeg&q=80&fit=cover
```

### Metadata

```
GET /api/image?url=<image_url>&metadata=true
```

Example Response:

```json
{
"width": 800,
"height": 600,
"format": "jpeg",
"mimeType": "image/jpeg"
}
```

### Placeholder

```
GET /api/image?url=<image_url>&placeholder=true&w=<width>&h=<height>&q=<quality>
```

### Structure

```
.
├── main.go # Server and handler implementation
├── internal/
│ ├── cache/ # Caching implementation
│ ├── devserver/ # Development server utilities
│ ├── metadata/ # Image metadata handling
│ ├── transform/ # Image transformation logic
│ ├── validate/ # Input validation
│ └── testdata/ # Test files and examples
└── .air.toml # Air configuration for hot reloading
```

### Running Tests

```bash
# Run all tests
CGO_CFLAGS="-Wno-xor-used-as-pow" go test -v ./...
# Run tests for specific package
CGO_CFLAGS="-Wno-xor-used-as-pow" go test -v ./internal/transform
go test -v ./internal/metadata
go test -v ./internal/validate
go test -v ./internal/cache
```

### Running the Server

```bash
GO_ENV=development go run main.go
```
