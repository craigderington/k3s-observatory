package k8s

import (
	"log"
	"math"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

// EventType represents the type of Kubernetes event
type EventType string

const (
	EventPodAdded    EventType = "pod_added"
	EventPodModified EventType = "pod_modified"
	EventPodDeleted  EventType = "pod_deleted"
	EventNodeAdded   EventType = "node_added"
	EventNodeModified EventType = "node_modified"
	EventNodeDeleted EventType = "node_deleted"
)

// WatchEvent represents a Kubernetes watch event
type WatchEvent struct {
	Type EventType   `json:"type"`
	Pod  *Pod        `json:"pod,omitempty"`
	Node *Node       `json:"node,omitempty"`
}

// WatchPods watches for pod changes and sends events to the channel
func (c *Client) WatchPods(events chan<- WatchEvent) error {
	watcher, err := c.Clientset.CoreV1().Pods("").Watch(c.ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	log.Println("Started watching pods...")

	go func() {
		defer watcher.Stop()

		for event := range watcher.ResultChan() {
			pod, ok := event.Object.(*corev1.Pod)
			if !ok {
				continue
			}

			var eventType EventType
			switch event.Type {
			case watch.Added:
				eventType = EventPodAdded
			case watch.Modified:
				eventType = EventPodModified
			case watch.Deleted:
				eventType = EventPodDeleted
			default:
				continue
			}

			// Convert to our Pod type
			simplePod := c.convertPod(pod)

			events <- WatchEvent{
				Type: eventType,
				Pod:  &simplePod,
			}

			log.Printf("Pod event: %s - %s/%s", eventType, pod.Namespace, pod.Name)
		}
	}()

	return nil
}

// WatchNodes watches for node changes and sends events to the channel
func (c *Client) WatchNodes(events chan<- WatchEvent) error {
	watcher, err := c.Clientset.CoreV1().Nodes().Watch(c.ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	log.Println("Started watching nodes...")

	go func() {
		defer watcher.Stop()

		for event := range watcher.ResultChan() {
			node, ok := event.Object.(*corev1.Node)
			if !ok {
				continue
			}

			var eventType EventType
			switch event.Type {
			case watch.Added:
				eventType = EventNodeAdded
			case watch.Modified:
				eventType = EventNodeModified
			case watch.Deleted:
				eventType = EventNodeDeleted
			default:
				continue
			}

			// Convert to our Node type
			simpleNode := c.convertNode(node, 0, 0)

			events <- WatchEvent{
				Type: eventType,
				Node: &simpleNode,
			}

			log.Printf("Node event: %s - %s", eventType, node.Name)
		}
	}()

	return nil
}

// Helper function to convert Kubernetes pod to our Pod type
func (c *Client) convertPod(kubePod *corev1.Pod) Pod {
	containers := make([]Container, 0, len(kubePod.Status.ContainerStatuses))
	for _, cs := range kubePod.Status.ContainerStatuses {
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

	// Simple position calculation (will be improved later)
	angle := 0.0
	orbitRadius := 3.0

	return Pod{
		ID:         string(kubePod.UID),
		Name:       kubePod.Name,
		Namespace:  kubePod.Namespace,
		Status:     string(kubePod.Status.Phase),
		NodeName:   kubePod.Spec.NodeName,
		Containers: containers,
		CreatedAt:  kubePod.CreationTimestamp.Time,
		Position: Position{
			X: orbitRadius * math.Cos(angle),
			Y: 0,
			Z: orbitRadius * math.Sin(angle),
		},
	}
}

// Helper function to convert Kubernetes node to our Node type
func (c *Client) convertNode(kubeNode *corev1.Node, index, total int) Node {
	// Determine node status
	status := "NotReady"
	for _, condition := range kubeNode.Status.Conditions {
		if condition.Type == "Ready" && condition.Status == "True" {
			status = "Ready"
			break
		}
	}

	// Get resource capacity
	cpu := kubeNode.Status.Capacity.Cpu().AsApproximateFloat64()
	memory := kubeNode.Status.Capacity.Memory().AsApproximateFloat64() / (1024 * 1024 * 1024)

	// Calculate position
	angle := 0.0
	if total > 0 {
		angle = float64(index) * 2.0 * math.Pi / float64(total)
	}
	radius := 10.0

	return Node{
		ID:     string(kubeNode.UID),
		Name:   kubeNode.Name,
		Status: status,
		CPU: ResourceUsage{
			Used:  0,
			Total: cpu,
		},
		Memory: ResourceUsage{
			Used:  0,
			Total: memory,
		},
		Pods:   []string{},
		Labels: kubeNode.Labels,
		Position: Position{
			X: radius * math.Cos(angle),
			Y: 0,
			Z: radius * math.Sin(angle),
		},
	}
}
