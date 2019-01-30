package main

import (
	"fmt"
	"log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/config"
)

const (
	kindImage   = "kindest/node:v1.13.2"
	clusterName = "godays"
)

// createCluster uses KinD to create a new Kubernetes cluster
// The cluster name is automatically prefixed with the 'kind-' prefix by NewContext function
func createCluster(name string, workerNodes int32) (*cluster.Context, error) {
	// Step 1: Create a new KinD cluster management context
	ctx := cluster.NewContext(name)

	// Step 2: Create a cluster config object describing how the cluster should look like
	cfg := &config.Config{
		Nodes: []config.Node{
			// Control plane node
			{
				// HA support is work-in-progress, so we use only one replica
				Replicas: int32ptr(1),
				Role:     config.ControlPlaneRole,
				Image:    kindImage,
			},
			// Worker nodes
			{
				Replicas: int32ptr(workerNodes),
				Role:     config.WorkerRole,
				Image:    kindImage,
			},
		},
	}

	// Step 3: Start the cluster
	// If retain (second argument) is true, nodes will not be destroyed in case of failure
	err := ctx.Create(cfg, false, 0)
	if err != nil {
		return nil, err
	}

	return ctx, nil
}

// getKubernetesClientset returns clientset used to communicate with the cluster
func getKubernetesClientset(ctx *cluster.Context) (*clientset.Clientset, error) {
	// Step 1: Parse Kubeconfig file using Kubernetes client-go
	cfg, err := clientcmd.BuildConfigFromFlags("", ctx.KubeConfigPath())
	if err != nil {
		return nil, err
	}

	// Step 2: Create a new Clientset so we can interact with the cluster
	return clientset.NewForConfig(cfg)
}

// deleteCluster deletes a cluster with the provided name
func deleteCluster(name string) error {
	// Step 1: Create a new KinD cluster management context
	ctx := cluster.NewContext(name)

	// Step 2: Delete the cluster
	return ctx.Delete()
}

func main() {
	// Create a new cluster
	c, err := createCluster(clusterName, 1)
	if err != nil {
		log.Fatal(err)
	}

	// Get Kubernetes Clientset
	cs, err := getKubernetesClientset(c)
	if err != nil {
		log.Fatal(err)
	}

	// Get all pods in the 'kube-system' namespace
	podList, err := cs.CoreV1().Pods("kube-system").List(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}

	// Go through pods list and append pod's name to a slice, then print all names
	var pods []string
	for _, pod := range podList.Items {
		pods = append(pods, pod.Name)
	}
	fmt.Println(pods)

	// Delete the cluster
	//err := deleteCluster(clusterName)
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func int32ptr(n int32) *int32 {
	return &n
}
