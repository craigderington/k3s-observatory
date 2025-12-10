package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	// Start metrics fetcher (5 second interval)
	metricsFetcher := k8s.NewMetricsFetcher(client, 5*time.Second)
	metricsChannel := metricsFetcher.Start()

	// Forward metrics updates to WebSocket clients
	go func() {
		for update := range metricsChannel {
			if err := hub.BroadcastEvent("metrics_update", update); err != nil {
				log.Printf("Error broadcasting metrics: %v", err)
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
	mux.HandleFunc("/api/pods/describe", apiHandler.DescribePod)
	mux.HandleFunc("/api/nodes/describe", apiHandler.DescribeNode)
	mux.HandleFunc("/api/pods/logs", apiHandler.GetPodLogs)
	mux.HandleFunc("/api/pods/metrics", apiHandler.GetPodMetrics)
	mux.HandleFunc("/api/nodes/metrics", apiHandler.GetNodeMetrics)
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.ServeWs(hub, w, r)
	})

	// Enable CORS for development
	corsHandler := enableCORS(mux)

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Shutting down gracefully...")
		metricsFetcher.Stop()
		os.Exit(0)
	}()

	// Start server
	port := getEnv("PORT", "8000")
	log.Printf("Observatory backend starting on port %s...", port)
	log.Printf("Endpoints available:")
	log.Printf("  GET /api/health")
	log.Printf("  GET /api/nodes")
	log.Printf("  GET /api/pods")
	log.Printf("  GET /api/pods/describe?namespace=X&name=Y")
	log.Printf("  GET /api/nodes/describe?name=X")
	log.Printf("  GET /api/pods/logs?namespace=X&name=Y&container=Z")
	log.Printf("  GET /api/pods/metrics?namespace=X&name=Y")
	log.Printf("  GET /api/nodes/metrics?name=X")
	log.Printf("  WS  /ws")
	log.Printf("Metrics fetcher running (5s interval)")

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
