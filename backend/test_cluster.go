package main

import (
	"context"
	"fmt"
	"log"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"os"
)

func main() {
	// Build config from default kubeconfig
	kubeconfig := os.Getenv("HOME") + "/.kube/config"
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		log.Fatal("Failed to build config:", err)
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal("Failed to create clientset:", err)
	}

	// Try to list pods
	fmt.Println("Attempting to list pods...")
	pods, err := clientset.CoreV1().Pods("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to list pods:", err)
	}
	fmt.Printf("Found %d pods\n", len(pods.Items))

	// Try to list nodes
	fmt.Println("Attempting to list nodes...")
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Fatal("Failed to list nodes:", err)
	}
	fmt.Printf("Found %d nodes\n", len(nodes.Items))
	
	fmt.Println("Cluster appears healthy!")
}
