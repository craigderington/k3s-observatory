package k8s

import (
	"log"
	"math"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GetPods fetches all pods from the cluster
func (c *Client) GetPods() ([]Pod, error) {
	podList, err := c.Clientset.CoreV1().Pods("").List(c.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// First get nodes to know their positions
	nodes, err := c.GetNodes()
	if err != nil {
		return nil, err
	}

	// Create a map of node names to their positions
	nodePositions := make(map[string]Position)
	for _, node := range nodes {
		nodePositions[node.Name] = node.Position
	}

	pods := make([]Pod, 0, len(podList.Items))

	// Group pods by node for positioning
	podsByNode := make(map[string]int)

	for _, pod := range podList.Items {
		// Extract container info
		containers := make([]Container, 0, len(pod.Status.ContainerStatuses))
		for _, cs := range pod.Status.ContainerStatuses {
			status := "Unknown"
			if cs.State.Running != nil {
				status = "Running"
			} else if cs.State.Waiting != nil {
				status = "Waiting"
			} else if cs.State.Terminated != nil {
				status = "Terminated"
			}

			containers = append(containers, Container{
				Name:     cs.Name,
				Status:   status,
				Restarts: cs.RestartCount,
			})
		}

		// Calculate position around the node
		nodeName := pod.Spec.NodeName
		podIndex := podsByNode[nodeName]
		totalPodsOnNode := 0

		// Count total pods on this node for better angle distribution
		for _, p := range podList.Items {
			if p.Spec.NodeName == nodeName {
				totalPodsOnNode++
			}
		}

		// Calculate orbit angle
		angle := float64(podIndex) * 2.0 * math.Pi / math.Max(float64(totalPodsOnNode), 1.0)
		orbitRadius := 3.0

		// Get node position (default to origin if not found)
		nodePos := Position{X: 0, Y: 0, Z: 0}
		if pos, ok := nodePositions[nodeName]; ok {
			nodePos = pos
		}

		// Position pod relative to its node
		pods = append(pods, Pod{
			ID:         string(pod.UID),
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Status:     string(pod.Status.Phase),
			NodeName:   nodeName,
			Containers: containers,
			CreatedAt:  pod.CreationTimestamp.Time,
			Position: Position{
				X: nodePos.X + orbitRadius*math.Cos(angle),
				Y: nodePos.Y,
				Z: nodePos.Z + orbitRadius*math.Sin(angle),
			},
		})

		podsByNode[nodeName]++
	}

	log.Printf("Fetched %d pods from cluster", len(pods))
	return pods, nil
}
