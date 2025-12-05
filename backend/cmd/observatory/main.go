package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/craigderington/lantern/internal/api"
	"github.com/craigderington/lantern/internal/k8s"
	"github.com/craigderington/lantern/internal/websocket"
)

func main() {
	// Initialize Kubernetes client
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	// Create WebSocket hub
	hub := websocket.NewHub()
	go hub.Run()

	// Start watching Kubernetes events
	events := make(chan k8s.WatchEvent, 256)

	if err := client.WatchPods(events); err != nil {
		log.Fatalf("Failed to start pod watcher: %v", err)
	}

	if err := client.WatchNodes(events); err != nil {
		log.Fatalf("Failed to start node watcher: %v", err)
	}

	// Forward Kubernetes events to WebSocket clients
	go func() {
		for event := range events {
			if err := hub.BroadcastEvent(string(event.Type), event); err != nil {
				log.Printf("Error broadcasting event: %v", err)
			}
		}
	}()

	// Create API handler
	apiHandler := api.NewHandler(client)

	// Setup routes
	mux := http.NewServeMux()
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("/api/nodes", apiHandler.GetNodes)
	mux.HandleFunc("/api/pods", apiHandler.GetPods)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Enable CORS for development
	corsHandler := enableCORS(mux)

	// Start server
	port := getEnv("PORT", "8000")
	log.Printf("Observatory backend starting on port %s...", port)
	log.Printf("Endpoints available:")
	log.Printf("  GET /api/health")
	log.Printf("  GET /api/nodes")
	log.Printf("  GET /api/pods")
	log.Printf("  WS  /ws")

	if err := http.ListenAndServe(":"+port, corsHandler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"service": "observatory-backend",
	})
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
