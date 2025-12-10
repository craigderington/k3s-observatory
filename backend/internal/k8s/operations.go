package k8s

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DescribePod returns a detailed description of a pod
func (c *Client) DescribePod(namespace, name string) (string, error) {
	pod, err := c.Clientset.CoreV1().Pods(namespace).Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get pod: %w", err)
	}

	var buf bytes.Buffer

	// Basic Info
	buf.WriteString(fmt.Sprintf("Name:         %s\n", pod.Name))
	buf.WriteString(fmt.Sprintf("Namespace:    %s\n", pod.Namespace))
	buf.WriteString(fmt.Sprintf("Node:         %s\n", pod.Spec.NodeName))
	buf.WriteString(fmt.Sprintf("Status:       %s\n", pod.Status.Phase))
	buf.WriteString(fmt.Sprintf("IP:           %s\n", pod.Status.PodIP))
	buf.WriteString(fmt.Sprintf("Created:      %s\n", pod.CreationTimestamp.Format("2006-01-02 15:04:05")))

	// Labels
	if len(pod.Labels) > 0 {
		buf.WriteString("\nLabels:\n")
		for k, v := range pod.Labels {
			buf.WriteString(fmt.Sprintf("  %s=%s\n", k, v))
		}
	}

	// Annotations
	if len(pod.Annotations) > 0 {
		buf.WriteString("\nAnnotations:\n")
		for k, v := range pod.Annotations {
			buf.WriteString(fmt.Sprintf("  %s=%s\n", k, v))
		}
	}

	// Containers
	buf.WriteString("\nContainers:\n")
	for _, container := range pod.Spec.Containers {
		buf.WriteString(fmt.Sprintf("  %s:\n", container.Name))
		buf.WriteString(fmt.Sprintf("    Image:         %s\n", container.Image))
		if len(container.Ports) > 0 {
			buf.WriteString("    Ports:\n")
			for _, port := range container.Ports {
				buf.WriteString(fmt.Sprintf("      %s/%d\n", port.Protocol, port.ContainerPort))
			}
		}
		if len(container.Env) > 0 {
			buf.WriteString("    Environment:\n")
			for _, env := range container.Env {
				buf.WriteString(fmt.Sprintf("      %s=%s\n", env.Name, env.Value))
			}
		}
	}

	// Container Statuses
	if len(pod.Status.ContainerStatuses) > 0 {
		buf.WriteString("\nContainer Statuses:\n")
		for _, cs := range pod.Status.ContainerStatuses {
			buf.WriteString(fmt.Sprintf("  %s:\n", cs.Name))
			buf.WriteString(fmt.Sprintf("    Ready:          %t\n", cs.Ready))
			buf.WriteString(fmt.Sprintf("    Restart Count:  %d\n", cs.RestartCount))
			if cs.State.Running != nil {
				buf.WriteString(fmt.Sprintf("    State:          Running (since %s)\n", cs.State.Running.StartedAt.Format("2006-01-02 15:04:05")))
			} else if cs.State.Waiting != nil {
				buf.WriteString(fmt.Sprintf("    State:          Waiting (%s)\n", cs.State.Waiting.Reason))
			} else if cs.State.Terminated != nil {
				buf.WriteString(fmt.Sprintf("    State:          Terminated (%s)\n", cs.State.Terminated.Reason))
			}
		}
	}

	// Conditions
	if len(pod.Status.Conditions) > 0 {
		buf.WriteString("\nConditions:\n")
		for _, cond := range pod.Status.Conditions {
			buf.WriteString(fmt.Sprintf("  %s:  %s\n", cond.Type, cond.Status))
		}
	}

	// Events
	events, err := c.Clientset.CoreV1().Events(namespace).List(c.ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("involvedObject.name=%s,involvedObject.namespace=%s", name, namespace),
	})
	if err == nil && len(events.Items) > 0 {
		buf.WriteString("\nRecent Events:\n")
		// Show last 10 events
		count := len(events.Items)
		if count > 10 {
			count = 10
		}
		for i := len(events.Items) - count; i < len(events.Items); i++ {
			event := events.Items[i]
			buf.WriteString(fmt.Sprintf("  %s  %s  %s\n",
				event.LastTimestamp.Format("15:04:05"),
				event.Reason,
				event.Message))
		}
	}

	return buf.String(), nil
}

// DescribeNode returns a detailed description of a node
func (c *Client) DescribeNode(name string) (string, error) {
	node, err := c.Clientset.CoreV1().Nodes().Get(c.ctx, name, metav1.GetOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to get node: %w", err)
	}

	var buf bytes.Buffer

	// Basic Info
	buf.WriteString(fmt.Sprintf("Name:               %s\n", node.Name))
	buf.WriteString(fmt.Sprintf("Created:            %s\n", node.CreationTimestamp.Format("2006-01-02 15:04:05")))

	// Labels
	if len(node.Labels) > 0 {
		buf.WriteString("\nLabels:\n")
		for k, v := range node.Labels {
			buf.WriteString(fmt.Sprintf("  %s=%s\n", k, v))
		}
	}

	// Conditions
	if len(node.Status.Conditions) > 0 {
		buf.WriteString("\nConditions:\n")
		for _, cond := range node.Status.Conditions {
			buf.WriteString(fmt.Sprintf("  %s:  %s  (%s)\n", cond.Type, cond.Status, cond.Reason))
		}
	}

	// Addresses
	if len(node.Status.Addresses) > 0 {
		buf.WriteString("\nAddresses:\n")
		for _, addr := range node.Status.Addresses {
			buf.WriteString(fmt.Sprintf("  %s:  %s\n", addr.Type, addr.Address))
		}
	}

	// Capacity
	buf.WriteString("\nCapacity:\n")
	buf.WriteString(fmt.Sprintf("  CPU:     %s\n", node.Status.Capacity.Cpu().String()))
	buf.WriteString(fmt.Sprintf("  Memory:  %s\n", node.Status.Capacity.Memory().String()))
	buf.WriteString(fmt.Sprintf("  Pods:    %s\n", node.Status.Capacity.Pods().String()))

	// Allocatable
	buf.WriteString("\nAllocatable:\n")
	buf.WriteString(fmt.Sprintf("  CPU:     %s\n", node.Status.Allocatable.Cpu().String()))
	buf.WriteString(fmt.Sprintf("  Memory:  %s\n", node.Status.Allocatable.Memory().String()))
	buf.WriteString(fmt.Sprintf("  Pods:    %s\n", node.Status.Allocatable.Pods().String()))

	// System Info
	buf.WriteString("\nSystem Info:\n")
	buf.WriteString(fmt.Sprintf("  OS:                   %s\n", node.Status.NodeInfo.OSImage))
	buf.WriteString(fmt.Sprintf("  Kernel:               %s\n", node.Status.NodeInfo.KernelVersion))
	buf.WriteString(fmt.Sprintf("  Container Runtime:    %s\n", node.Status.NodeInfo.ContainerRuntimeVersion))
	buf.WriteString(fmt.Sprintf("  Kubelet Version:      %s\n", node.Status.NodeInfo.KubeletVersion))

	// List pods on this node
	pods, err := c.Clientset.CoreV1().Pods("").List(c.ctx, metav1.ListOptions{
		FieldSelector: fmt.Sprintf("spec.nodeName=%s", name),
	})
	if err == nil {
		buf.WriteString(fmt.Sprintf("\nPods: (%d)\n", len(pods.Items)))
		for _, pod := range pods.Items {
			buf.WriteString(fmt.Sprintf("  %s/%s  (%s)\n", pod.Namespace, pod.Name, pod.Status.Phase))
		}
	}

	return buf.String(), nil
}

// GetPodLogs returns logs for a specific pod container
func (c *Client) GetPodLogs(namespace, podName, containerName string, tailLines int64) (string, error) {
	opts := &corev1.PodLogOptions{
		Container: containerName,
		TailLines: &tailLines,
	}

	req := c.Clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	podLogs, err := req.Stream(context.Background())
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	defer podLogs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}

	return buf.String(), nil
}

// PodMetrics represents resource metrics for a pod
type PodMetrics struct {
	Name        string  `json:"name"`
	Namespace   string  `json:"namespace"`
	CPUUsage    float64 `json:"cpuUsage"`    // in millicores
	MemoryUsage float64 `json:"memoryUsage"` // in MB
	Timestamp   string  `json:"timestamp"`
}

// NodeMetrics represents resource metrics for a node
type NodeMetrics struct {
	Name        string  `json:"name"`
	CPUUsage    float64 `json:"cpuUsage"`    // in millicores
	MemoryUsage float64 `json:"memoryUsage"` // in MB
	Timestamp   string  `json:"timestamp"`
}

// GetPodMetrics fetches current resource usage for a pod
func (c *Client) GetPodMetrics(namespace, name string) (*PodMetrics, error) {
	// Use the metrics.k8s.io API
	data, err := c.Clientset.RESTClient().
		Get().
		AbsPath("/apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods/" + name).
		DoRaw(c.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	// Parse the metrics response
	var result struct {
		Metadata struct {
			Name      string `json:"name"`
			Namespace string `json:"namespace"`
		} `json:"metadata"`
		Timestamp  string `json:"timestamp"`
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	// Aggregate container metrics
	var totalCPU, totalMemory float64
	for _, container := range result.Containers {
		// Parse CPU (format: "123n" or "456m")
		cpuStr := container.Usage.CPU
		if strings.HasSuffix(cpuStr, "n") {
			// nanocores to millicores
			val, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "n"), 64)
			totalCPU += val / 1000000
		} else if strings.HasSuffix(cpuStr, "m") {
			// millicores
			val, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "m"), 64)
			totalCPU += val
		}

		// Parse Memory (format: "123456Ki" or "789Mi")
		memStr := container.Usage.Memory
		if strings.HasSuffix(memStr, "Ki") {
			// KiB to MB
			val, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Ki"), 64)
			totalMemory += val / 1024
		} else if strings.HasSuffix(memStr, "Mi") {
			// MiB to MB (approximately)
			val, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Mi"), 64)
			totalMemory += val
		}
	}

	return &PodMetrics{
		Name:        result.Metadata.Name,
		Namespace:   result.Metadata.Namespace,
		CPUUsage:    totalCPU,
		MemoryUsage: totalMemory,
		Timestamp:   result.Timestamp,
	}, nil
}

// GetNodeMetrics fetches current resource usage for a node
func (c *Client) GetNodeMetrics(name string) (*NodeMetrics, error) {
	data, err := c.Clientset.RESTClient().
		Get().
		AbsPath("/apis/metrics.k8s.io/v1beta1/nodes/" + name).
		DoRaw(c.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get node metrics: %w", err)
	}

	var result struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Timestamp string `json:"timestamp"`
		Usage     struct {
			CPU    string `json:"cpu"`
			Memory string `json:"memory"`
		} `json:"usage"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	// Parse CPU
	var cpuUsage float64
	cpuStr := result.Usage.CPU
	if strings.HasSuffix(cpuStr, "n") {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "n"), 64)
		cpuUsage = val / 1000000
	} else if strings.HasSuffix(cpuStr, "m") {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(cpuStr, "m"), 64)
		cpuUsage = val
	}

	// Parse Memory
	var memoryUsage float64
	memStr := result.Usage.Memory
	if strings.HasSuffix(memStr, "Ki") {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Ki"), 64)
		memoryUsage = val / 1024
	} else if strings.HasSuffix(memStr, "Mi") {
		val, _ := strconv.ParseFloat(strings.TrimSuffix(memStr, "Mi"), 64)
		memoryUsage = val
	}

	return &NodeMetrics{
		Name:        result.Metadata.Name,
		CPUUsage:    cpuUsage,
		MemoryUsage: memoryUsage,
		Timestamp:   result.Timestamp,
	}, nil
}

// ContainerMetricsData represents metrics for a single container
type ContainerMetricsData struct {
	Name   string  `json:"name"`
	CPU    float64 `json:"cpu"`    // millicores
	Memory float64 `json:"memory"` // MB
}

// GetPodMetricsDetailed fetches per-container resource usage for a pod
func (c *Client) GetPodMetricsDetailed(namespace, name string) ([]ContainerMetricsData, error) {
	// Use the metrics.k8s.io API
	data, err := c.Clientset.RESTClient().
		Get().
		AbsPath("/apis/metrics.k8s.io/v1beta1/namespaces/" + namespace + "/pods/" + name).
		DoRaw(c.ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get pod metrics: %w", err)
	}

	// Parse the metrics response
	var result struct {
		Containers []struct {
			Name  string `json:"name"`
			Usage struct {
				CPU    string `json:"cpu"`
				Memory string `json:"memory"`
			} `json:"usage"`
		} `json:"containers"`
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse metrics: %w", err)
	}

	// Convert to our type
	containerMetrics := make([]ContainerMetricsData, 0, len(result.Containers))
	for _, container := range result.Containers {
		cpu := parseMetricValue(container.Usage.CPU, "cpu")
		memory := parseMetricValue(container.Usage.Memory, "memory")

		containerMetrics = append(containerMetrics, ContainerMetricsData{
			Name:   container.Name,
			CPU:    cpu,
			Memory: memory,
		})
	}

	return containerMetrics, nil
}

// parseMetricValue parses Kubernetes metric values into float64
func parseMetricValue(value string, metricType string) float64 {
	if metricType == "cpu" {
		if strings.HasSuffix(value, "n") {
			// nanocores to millicores
			val, _ := strconv.ParseFloat(strings.TrimSuffix(value, "n"), 64)
			return val / 1000000
		} else if strings.HasSuffix(value, "m") {
			// millicores
			val, _ := strconv.ParseFloat(strings.TrimSuffix(value, "m"), 64)
			return val
		}
	} else if metricType == "memory" {
		if strings.HasSuffix(value, "Ki") {
			// KiB to MB
			val, _ := strconv.ParseFloat(strings.TrimSuffix(value, "Ki"), 64)
			return val / 1024
		} else if strings.HasSuffix(value, "Mi") {
			// MiB
			val, _ := strconv.ParseFloat(strings.TrimSuffix(value, "Mi"), 64)
			return val
		}
	}
	return 0
}
