package devserver

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = struct {
	sync.RWMutex
	connections []*websocket.Conn
}{}

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

	// Add WebSocket endpoint
	mux.HandleFunc("/ws", handleWebSocket)

	// Start watching for reload triggers
	go watchReloadTrigger()

	log.Printf("Development mode: Test page available at http://localhost:%s/test/test.html", port)
	return nil
}

func watchReloadTrigger() {
	triggerFile := "tmp/reload-trigger"
	var lastMod time.Time

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	log.Printf("Starting reload trigger watcher")

	for range ticker.C {
		stat, err := os.Stat(triggerFile)
		if err != nil {
			continue // File doesn't exist yet
		}

		modTime := stat.ModTime()
		if modTime != lastMod {
			lastMod = modTime
			log.Printf("Reload trigger detected")
			notifyClients()
		}
	}
}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	log.Printf("New WebSocket connection established")

	clients.Lock()
	clients.connections = append(clients.connections, conn)
	clients.Unlock()

	// Keep connection alive and remove it when closed
	go func() {
		defer func() {
			conn.Close()
			clients.Lock()
			for i, c := range clients.connections {
				if c == conn {
					clients.connections = append(clients.connections[:i], clients.connections[i+1:]...)
					break
				}
			}
			clients.Unlock()
			log.Printf("WebSocket connection closed")
		}()

		for {
			// Read messages to keep connection alive
			if _, _, err := conn.ReadMessage(); err != nil {
				return
			}
		}
	}()
}

func notifyClients() {
	clients.Lock()
	defer clients.Unlock()

	log.Printf("Notifying %d clients to reload", len(clients.connections))
	for i := len(clients.connections) - 1; i >= 0; i-- {
		if err := clients.connections[i].WriteMessage(
			websocket.TextMessage,
			[]byte("reload")); err != nil {
			// Remove dead connections
			clients.connections = append(
				clients.connections[:i],
				clients.connections[i+1:]...)
		}
	}
}