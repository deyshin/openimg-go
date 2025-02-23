package main

import (
	"fmt"
	"image"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/yourusername/openimg-go/internal/transform"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create a new image handler
	handler := &ImageHandler{
		Client: &http.Client{},
	}

	// Register routes
	http.HandleFunc("/api/image", handler.ServeImage)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

type ImageHandler struct {
	Client *http.Client
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

	// Get image URL and transformation parameters
	imageURL := r.URL.Query().Get("url")
	if imageURL == "" {
		http.Error(w, "Missing url parameter", http.StatusBadRequest)
		return
	}

	// Parse transformation options
	width, _ := strconv.Atoi(r.URL.Query().Get("w"))
	height, _ := strconv.Atoi(r.URL.Query().Get("h"))
	quality, _ := strconv.Atoi(r.URL.Query().Get("q"))
	format := r.URL.Query().Get("fmt")
	fit := r.URL.Query().Get("fit")

	// Validate parameters
	if quality < 1 || quality > 100 {
		quality = 85 // default quality
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
	contentType := "image/jpeg"
	if format == "png" {
		contentType = "image/png"
	}
	w.Header().Set("Content-Type", contentType)
	w.Write(transformed)
}