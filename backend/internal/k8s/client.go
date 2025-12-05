package k8s

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Client struct {
	Clientset *kubernetes.Clientset
	ctx       context.Context
}

// NewClient creates a new Kubernetes client
// It tries in-cluster config first, then falls back to kubeconfig
func NewClient() (*Client, error) {
	var config *rest.Config
	var err error

	// Try in-cluster config first (for when running inside k8s)
	config, err = rest.InClusterConfig()
	if err != nil {
		// Fall back to kubeconfig (for local development)
		log.Println("Not running in cluster, using kubeconfig...")
		config, err = getKubeConfig()
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("Using in-cluster configuration")
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	log.Println("Successfully connected to Kubernetes cluster")

	return &Client{
		Clientset: clientset,
		ctx:       context.Background(),
	}, nil
}

// getKubeConfig loads kubeconfig from KUBECONFIG env var or default location
func getKubeConfig() (*rest.Config, error) {
	var kubeconfig string

	// Check KUBECONFIG environment variable first
	if envKubeconfig := os.Getenv("KUBECONFIG"); envKubeconfig != "" {
		kubeconfig = envKubeconfig
		log.Printf("Using KUBECONFIG from environment: %s", kubeconfig)
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
		log.Printf("Using default kubeconfig: %s", kubeconfig)
	}

	// Build config from kubeconfig file
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	return config, nil
}

// Context returns the client's context
func (c *Client) Context() context.Context {
	return c.ctx
}
