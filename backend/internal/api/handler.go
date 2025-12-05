package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/craigderington/lantern/internal/k8s"
)

type Handler struct {
	k8sClient *k8s.Client
}

func NewHandler(client *k8s.Client) *Handler {
	return &Handler{
		k8sClient: client,
	}
}

// GetNodes handles GET /api/nodes
func (h *Handler) GetNodes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nodes, err := h.k8sClient.GetNodes()
	if err != nil {
		log.Printf("Error fetching nodes: %v", err)
		http.Error(w, "Failed to fetch nodes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(nodes); err != nil {
		log.Printf("Error encoding nodes response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Served %d nodes", len(nodes))
}

// GetPods handles GET /api/pods
func (h *Handler) GetPods(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pods, err := h.k8sClient.GetPods()
	if err != nil {
		log.Printf("Error fetching pods: %v", err)
		http.Error(w, "Failed to fetch pods", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pods); err != nil {
		log.Printf("Error encoding pods response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Served %d pods", len(pods))
}
