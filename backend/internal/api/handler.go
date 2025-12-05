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

// DescribePod handles GET /api/pods/{namespace}/{name}/describe
func (h *Handler) DescribePod(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract namespace and name from URL path
	// Expected format: /api/pods/{namespace}/{name}/describe
	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" || name == "" {
		http.Error(w, "namespace and name query parameters required", http.StatusBadRequest)
		return
	}

	description, err := h.k8sClient.DescribePod(namespace, name)
	if err != nil {
		log.Printf("Error describing pod %s/%s: %v", namespace, name, err)
		http.Error(w, "Failed to describe pod", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(description))

	log.Printf("Described pod: %s/%s", namespace, name)
}

// DescribeNode handles GET /api/nodes/{name}/describe
func (h *Handler) DescribeNode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name query parameter required", http.StatusBadRequest)
		return
	}

	description, err := h.k8sClient.DescribeNode(name)
	if err != nil {
		log.Printf("Error describing node %s: %v", name, err)
		http.Error(w, "Failed to describe node", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(description))

	log.Printf("Described node: %s", name)
}

// GetPodLogs handles GET /api/pods/{namespace}/{name}/logs
func (h *Handler) GetPodLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")
	container := r.URL.Query().Get("container")

	if namespace == "" || name == "" {
		http.Error(w, "namespace and name query parameters required", http.StatusBadRequest)
		return
	}

	// Default to 100 lines
	tailLines := int64(100)

	logs, err := h.k8sClient.GetPodLogs(namespace, name, container, tailLines)
	if err != nil {
		log.Printf("Error getting logs for pod %s/%s: %v", namespace, name, err)
		http.Error(w, "Failed to get pod logs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(logs))

	log.Printf("Served logs for pod: %s/%s (container: %s)", namespace, name, container)
}

// GetPodMetrics handles GET /api/pods/metrics
func (h *Handler) GetPodMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	namespace := r.URL.Query().Get("namespace")
	name := r.URL.Query().Get("name")

	if namespace == "" || name == "" {
		http.Error(w, "namespace and name query parameters required", http.StatusBadRequest)
		return
	}

	metrics, err := h.k8sClient.GetPodMetrics(namespace, name)
	if err != nil {
		log.Printf("Error getting metrics for pod %s/%s: %v", namespace, name, err)
		http.Error(w, "Failed to get pod metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Served metrics for pod: %s/%s", namespace, name)
}

// GetNodeMetrics handles GET /api/nodes/metrics
func (h *Handler) GetNodeMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "name query parameter required", http.StatusBadRequest)
		return
	}

	metrics, err := h.k8sClient.GetNodeMetrics(name)
	if err != nil {
		log.Printf("Error getting metrics for node %s: %v", name, err)
		http.Error(w, "Failed to get node metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Printf("Served metrics for node: %s", name)
}
