package k8s

import (
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetNodes fetches all nodes from the cluster
func (c *Client) GetNodes() ([]Node, error) {
	nodeList, err := c.Clientset.CoreV1().Nodes().List(c.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	nodes := make([]Node, 0, len(nodeList.Items))

	for i, node := range nodeList.Items {
		// Determine node status
		status := "NotReady"
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" && condition.Status == "True" {
				status = "Ready"
				break
			}
		}

		// Get resource capacity
		cpu := node.Status.Capacity.Cpu().AsApproximateFloat64()
		memory := node.Status.Capacity.Memory().AsApproximateFloat64() / (1024 * 1024 * 1024) // Convert to GB

		// Calculate position in a circle
		angle := float64(i) * 2.0 * 3.14159 / float64(len(nodeList.Items))
		radius := 10.0

		nodes = append(nodes, Node{
			ID:     string(node.UID),
			Name:   node.Name,
			Status: status,
			CPU: ResourceUsage{
				Used:  0, // Will be populated with metrics later
				Total: cpu,
			},
			Memory: ResourceUsage{
				Used:  0, // Will be populated with metrics later
				Total: memory,
			},
			Pods:   []string{}, // Will be populated when fetching pods
			Labels: node.Labels,
			Position: Position{
				X: radius * cos(angle),
				Y: 0,
				Z: radius * sin(angle),
			},
		})
	}

	log.Printf("Fetched %d nodes from cluster", len(nodes))
	return nodes, nil
}

// Simple math helpers
func cos(x float64) float64 {
	// Simple cosine approximation for small angles
	return 1.0 - (x*x)/2.0 + (x*x*x*x)/24.0
}

func sin(x float64) float64 {
	// Simple sine approximation
	return x - (x*x*x)/6.0 + (x*x*x*x*x)/120.0
}
