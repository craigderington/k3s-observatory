package k8s

import "time"

// Node represents a simplified Kubernetes node for the frontend
type Node struct {
	ID       string            `json:"id"`
	Name     string            `json:"name"`
	Status   string            `json:"status"`
	CPU      ResourceUsage     `json:"cpu"`
	Memory   ResourceUsage     `json:"memory"`
	Pods     []string          `json:"pods"`
	Labels   map[string]string `json:"labels"`
	Position Position          `json:"position"`
}

// Pod represents a simplified Kubernetes pod for the frontend
type Pod struct {
	ID         string      `json:"id"`
	Name       string      `json:"name"`
	Namespace  string      `json:"namespace"`
	Status     string      `json:"status"`
	NodeName   string      `json:"nodeName"`
	Containers []Container `json:"containers"`
	CreatedAt  time.Time   `json:"createdAt"`
	Position   Position    `json:"position"`
}

// Container represents a container within a pod
type Container struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Restarts int32  `json:"restarts"`
}

// ResourceUsage represents resource consumption
type ResourceUsage struct {
	Used  float64 `json:"used"`
	Total float64 `json:"total"`
}

// Position represents 3D coordinates
type Position struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	Z float64 `json:"z"`
}
