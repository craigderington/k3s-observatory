package k8s

import (
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodMetricsData represents metrics for a single pod
type PodMetricsData struct {
	PodID            string                   `json:"podId"`
	Name             string                   `json:"name"`
	Namespace        string                   `json:"namespace"`
	TotalCPU         float64                  `json:"totalCpu"`
	TotalMemory      float64                  `json:"totalMemory"`
	ContainerMetrics []ContainerMetricsData   `json:"containers"`
	Timestamp        time.Time                `json:"timestamp"`
}

// MetricsUpdate represents a batch metrics update
type MetricsUpdate struct {
	Type      string            `json:"type"`
	Pods      []PodMetricsData  `json:"pods"`
	Timestamp time.Time         `json:"timestamp"`
}

// MetricsFetcher periodically fetches and broadcasts pod metrics
type MetricsFetcher struct {
	client   *Client
	interval time.Duration
	stop     chan struct{}
}

// NewMetricsFetcher creates a new metrics fetcher
func NewMetricsFetcher(client *Client, interval time.Duration) *MetricsFetcher {
	return &MetricsFetcher{
		client:   client,
		interval: interval,
		stop:     make(chan struct{}),
	}
}

// Start begins the periodic metrics collection
func (mf *MetricsFetcher) Start() <-chan MetricsUpdate {
	updates := make(chan MetricsUpdate, 256)

	go func() {
		ticker := time.NewTicker(mf.interval)
		defer ticker.Stop()
		defer close(updates)

		for {
			select {
			case <-ticker.C:
				mf.fetchAndBroadcast(updates)
			case <-mf.stop:
				log.Println("Metrics fetcher stopped")
				return
			}
		}
	}()

	return updates
}

// Stop halts the metrics fetcher
func (mf *MetricsFetcher) Stop() {
	close(mf.stop)
}

// fetchAndBroadcast fetches metrics for all pods
func (mf *MetricsFetcher) fetchAndBroadcast(updates chan<- MetricsUpdate) {
	// Get all pods
	pods, err := mf.client.Clientset.CoreV1().Pods("").List(mf.client.ctx, metav1.ListOptions{})
	if err != nil {
		log.Printf("Error listing pods for metrics: %v", err)
		return
	}

	podMetrics := make([]PodMetricsData, 0, len(pods.Items))

	for _, pod := range pods.Items {
		// Fetch metrics for this pod
		metrics, err := mf.client.GetPodMetrics(pod.Namespace, pod.Name)
		if err != nil {
			// Metrics might not be available for all pods (pending, completed, etc.)
			// Log at debug level and skip
			continue
		}

		// Also get per-container metrics
		containerMetrics, err := mf.client.GetPodMetricsDetailed(pod.Namespace, pod.Name)
		if err != nil {
			// Skip if detailed metrics not available
			continue
		}

		podMetrics = append(podMetrics, PodMetricsData{
			PodID:            string(pod.UID),
			Name:             pod.Name,
			Namespace:        pod.Namespace,
			TotalCPU:         metrics.CPUUsage,
			TotalMemory:      metrics.MemoryUsage,
			ContainerMetrics: containerMetrics,
			Timestamp:        time.Now(),
		})
	}

	if len(podMetrics) > 0 {
		updates <- MetricsUpdate{
			Type:      "metrics_update",
			Pods:      podMetrics,
			Timestamp: time.Now(),
		}
		log.Printf("Broadcasted metrics for %d pods", len(podMetrics))
	}
}
