package devserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Setup configures development server routes and features
func Setup(mux *http.ServeMux, port string) error {
	// Get the absolute path to the testdata directory
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	testDataDir := filepath.Join(wd, "internal", "testdata")

	// Serve files from the testdata directory
	fs := http.FileServer(http.Dir(testDataDir))
	mux.Handle("/test/", http.StripPrefix("/test/", fs))

	log.Printf("Development mode: Test page available at http://localhost:%s/test/test.html", port)
	return nil
}