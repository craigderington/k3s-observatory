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
		podIndex := podsByNode[pod.Spec.NodeName]
		podsByNode[pod.Spec.NodeName]++

		angle := float64(podIndex) * 2.0 * math.Pi / math.Max(float64(podsByNode[pod.Spec.NodeName]), 3.0)
		orbitRadius := 3.0

		pods = append(pods, Pod{
			ID:         string(pod.UID),
			Name:       pod.Name,
			Namespace:  pod.Namespace,
			Status:     string(pod.Status.Phase),
			NodeName:   pod.Spec.NodeName,
			Containers: containers,
			CreatedAt:  pod.CreationTimestamp.Time,
			Position: Position{
				X: orbitRadius * math.Cos(angle),
				Y: 0,
				Z: orbitRadius * math.Sin(angle),
			},
		})
	}

	log.Printf("Fetched %d pods from cluster", len(pods))
	return pods, nil
}
