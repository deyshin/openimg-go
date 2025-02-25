package main

import (
	"encoding/json"
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/yourusername/openimg-go/internal/cache"
	"github.com/yourusername/openimg-go/internal/devserver"
	"github.com/yourusername/openimg-go/internal/metadata"
	"github.com/yourusername/openimg-go/internal/transform"
	"github.com/yourusername/openimg-go/internal/validate"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new image handler
	handler := &ImageHandler{
		Client: &http.Client{},
		Cache:  cache.New(),
	}

	// Register routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/image", handler.ServeImage)

	// In development mode, serve test files
	if os.Getenv("GO_ENV") != "production" {
		log.Printf("Initializing development mode...")
		if err := devserver.Setup(mux, port); err != nil {
			log.Fatal(err)
		}
		log.Printf("Development mode initialized")
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

type ImageHandler struct {
	Client *http.Client
	Cache  *cache.ImageCache
}

func (h *ImageHandler) ServeImage(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if metadata is requested
	if r.URL.Query().Get("metadata") == "true" {
		h.serveMetadata(w, r)
		return
	}

	// Check if placeholder is requested
	if r.URL.Query().Get("placeholder") == "true" {
		h.servePlaceholder(w, r)
		return
	}

	// Get image URL and transformation parameters
	imageURL := r.URL.Query().Get("url")
	if err := validate.URL(imageURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse transformation options
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	height, _ := strconv.Atoi(r.URL.Query().Get("h"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("q"))
	format := r.URL.Query().Get("fmt")
	fit := r.URL.Query().Get("fit")

	if err := validate.ImageOptions(width, height, quality, format, fit); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate parameters
	if quality < 1 || quality > 100 {
		quality = 85 // default quality
	}

	// Generate cache key
	cacheKey := cache.GenerateKey(imageURL, width, height, quality, format, fit)

	// Try to get from cache
	if cached, found := h.Cache.Get(cacheKey); found {
		contentType := "image/png"
		switch format {
		case "jpg", "jpeg":
			contentType = "image/jpeg"
		case "png":
		case "avif":
			contentType = "image/avif"
		case "webp":
			contentType = "image/webp"
		}
		w.Header().Set("Content-Type", contentType)
		w.Write(cached)
		return
	}

	// Fetch the image
	resp, err := h.Client.Get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// Decode the image
	img, imgFormat, err := image.Decode(resp.Body)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}

	// If format is not specified, use original format
	if format == "" {
		format = imgFormat
	}

	// Transform the image
	transformed, err := transform.Transform(img, transform.Options{
		Width:   width,
		Height:  height,
		Format:  format,
		Quality: quality,
		Fit:     fit,
	})
	if err != nil {
		http.Error(w, "Failed to transform image", http.StatusInternalServerError)
		return
	}

	// Set content type based on format
	contentType := "image/png" // default
	switch format {
	case "jpg", "jpeg":
		contentType = "image/jpeg"
	case "png":
	case "avif":
		contentType = "image/avif"
	case "webp":
		contentType = "image/webp"
	}

	// Store in cache
	h.Cache.Set(cacheKey, transformed)

	w.Header().Set("Content-Type", contentType)
	w.Write(transformed)
}

func (h *ImageHandler) serveMetadata(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("url")
	if err := validate.URL(imageURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.Client.Get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	meta, err := metadata.Get(resp.Body)
	if err != nil {
		http.Error(w, "Failed to get image metadata", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(meta)
}

func (h *ImageHandler) servePlaceholder(w http.ResponseWriter, r *http.Request) {
	imageURL := r.URL.Query().Get("url")
	if err := validate.URL(imageURL); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Parse placeholder options
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	height, _ := strconv.Atoi(r.URL.Query().Get("h"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("q"))

	// Generate cache key for placeholder
	cacheKey := cache.GenerateKey(imageURL, width, height, quality, "placeholder", "")

	// Try to get from cache
	if cached, found := h.Cache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(cached)
		return
	}

	// Fetch and decode the image
	resp, err := h.Client.Get(imageURL)
	if err != nil {
		http.Error(w, "Failed to fetch image", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		http.Error(w, "Failed to decode image", http.StatusBadRequest)
		return
	}

	// Generate placeholder
	placeholder, err := transform.GeneratePlaceholder(img, transform.PlaceholderOptions{
		Width:   width,
		Height:  height,
		Quality: quality,
	})
	if err != nil {
		http.Error(w, "Failed to generate placeholder", http.StatusInternalServerError)
		return
	}

	// Store in cache
	h.Cache.Set(cacheKey, []byte(placeholder))

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(placeholder))
}