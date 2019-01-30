package main

import (
	"log"

	"github.com/kubicorn/kubicorn/apis/cluster"
	"github.com/kubicorn/kubicorn/pkg"
	"github.com/kubicorn/kubicorn/pkg/initapi"
	"github.com/kubicorn/kubicorn/pkg/logger"
	"github.com/kubicorn/kubicorn/profiles/digitalocean"
)

const (
	clusterName = "godays"
)

// init initializes a Kubicorn logger
func init() {
	logger.Level = 4
	logger.Info("Init cluster...")
}

// createCluster creates a new Kubernetes cluster using Kubicorn
func createCluster(name string) (*cluster.Cluster, error) {
	// Step 1: Initialize a profile for a new DigitalOcean cluster
	profile := digitalocean.NewUbuntuCluster(name)

	// Optional: Customize the profile
	config := profile.ProviderConfig()
	config.SSH.PublicKeyPath = "~/.ssh/k8s_rsa.pub"
	err := profile.SetProviderConfig(config)
	if err != nil {
		return nil, err
	}

	// Step 2: Run preprocessors and validation
	profile, err = initapi.InitCluster(profile)
	if err != nil {
		return nil, err
	}

	// Step 3: Get reconciler used to create and manage a cluster
	reconciler, err := pkg.GetReconciler(profile, nil)
	if err != nil {
		return nil, err
	}

	// Step 4: Get expected state (what we want) and actual state (what we already have on cloud)
	expected, err := reconciler.Expected(profile)
	if err != nil {
		return nil, err
	}
	actual, err := reconciler.Actual(profile)
	if err != nil {
		return nil, err
	}

	// Step 5: Create (reconcile) the cluster
	return reconciler.Reconcile(actual, expected)
}

// deleteCluster deletes the cluster
func deleteCluster(name string) error {
	// Step 1: Initialize a profile for a DigitalOcean profile
	profile := digitalocean.NewUbuntuCluster(name)

	// Optional: Customize the profile
	config := profile.ProviderConfig()
	config.SSH.PublicKeyPath = "~/.ssh/k8s_rsa.pub"
	err := profile.SetProviderConfig(config)
	if err != nil {
		return err
	}

	// Step 2: Run preprocessors and validation
	profile, err = initapi.InitCluster(profile)
	if err != nil {
		return err
	}

	// Step 3: Get reconciler used to create and manage a cluster
	reconciler, err := pkg.GetReconciler(profile, nil)
	if err != nil {
		return err
	}

	// Step 4: Destroy the cluster
	_, err = reconciler.Destroy()
	return err
}

func main() {
	_, err := createCluster(clusterName)
	if err != nil {
		log.Fatal(err)
	}

	err = deleteCluster(clusterName)
	if err != nil {
		log.Fatal(err)
	}
}
